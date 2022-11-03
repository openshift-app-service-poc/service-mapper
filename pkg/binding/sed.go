package binding

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bindingoperatorscoreoscomv1alpha1 "github.com/filariow/sbo-1225/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServiceEndpointDefinition(ctx context.Context,
	client client.Client,
	sm *bindingoperatorscoreoscomv1alpha1.ServiceResourceMap,
	sp *bindingoperatorscoreoscomv1alpha1.ServiceProxy,
	obj interface{}) *corev1.Secret {

	secrets := extractSecrets(ctx, client, sm.Namespace, sm.Spec.ServiceMap, obj)

	sed := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sp.Name + "-sed",
			Namespace: sp.Namespace,
		},
		StringData: secrets,
	}
	return &sed
}

func extractSecrets(ctx context.Context, client client.Client, namespace string, rules map[string]string, obj interface{}) map[string]string {
	secrets := map[string]string{}
	l, _ := logr.FromContext(ctx)

	for k, v := range rules {
		l.Info("processing secret rule", "key", k, "value", v)
		ss, err := processTarget(ctx, client, namespace, k, v, obj)
		if err != nil {
			l.Error(err, "can not process target", "target", obj, "rule", v)
			continue
		}
		l.Info("processed secret rule", "key", k, "value", v, "secret", ss)

		for k, v := range ss {
			secrets[k] = v
		}
	}

	return secrets
}

func processTarget(ctx context.Context, client client.Client, namespace, k, v string, obj interface{}) (map[string]string, error) {
	ss := strings.Split(v, ",")
	if len(ss) > 1 {
		return processRefTarget(ctx, client, namespace, ss, obj)
	}

	if isJsonpath(v) {
		v, err := executeJsonpath(v, obj)
		if err != nil {
			return nil, err
		}
		return map[string]string{k: v}, nil
	}
	return map[string]string{k: v}, nil
}

func processRefTarget(ctx context.Context, cli client.Client, namespace string, ss []string, obj interface{}) (map[string]string, error) {
	jp, refType := ss[0], ss[1]

	refObj, err := executeJsonpath(jp, obj)
	if err != nil {
		return nil, err
	}

	l, _ := logr.FromContext(ctx)
	l.Info("processing ref target", "refType", refType, "obj", obj)

	switch refType {
	case "objectType=Secret":
		s := corev1.Secret{}
		skey := client.ObjectKey{Namespace: namespace, Name: refObj}
		if err := cli.Get(ctx, skey, &s); err != nil {
			return nil, fmt.Errorf("can not retrieve Secret '%s/%s': %w", namespace, refObj, err)
		}

		l.Info("retrieved ref secret", "secret", s, "stringData", s.StringData)

		d := make(map[string]string, len(s.Data))
		for k, v := range s.Data {
			v64, err := base64.StdEncoding.DecodeString(string(v))
			if err != nil {
				l.Error(err, "can not base64 decode secret entry", "key", k)
			}

			d[k] = string(v64)
		}

		return d, nil
	case "objectType=ConfigMap":
		cm := corev1.ConfigMap{}
		cmkey := client.ObjectKey{Namespace: namespace, Name: refObj}
		if err := cli.Get(ctx, cmkey, &cm); err != nil {
			return nil, fmt.Errorf("can not retrieve ConfigMap '%s/%s': %w", namespace, refObj, err)
		}

		return cm.Data, nil
	}

	return nil, fmt.Errorf("invalid objectType: %s", refType)
}

func isJsonpath(v string) bool {
	return strings.HasPrefix(v, "path=")
}

func executeJsonpath(v string, data interface{}) (string, error) {
	jp := jsonpath.New("")

	if err := jp.Parse(v); err != nil {
		return "", fmt.Errorf("invalid jsonpath '%s': %w", v, err)
	}

	buf := new(bytes.Buffer)
	if err := jp.Execute(buf, data); err != nil {
		return "", fmt.Errorf("can not extract data using jsonpath '%s': %w", v, err)
	}

	return strings.TrimPrefix(buf.String(), "path="), nil
}
