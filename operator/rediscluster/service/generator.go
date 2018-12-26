package service

import (
	"fmt"

	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func generateRedisAccessService(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Service {
	name := generateName(svcAccessNamePrefix, rc.Name)
	namespace := rc.Namespace

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       redisAccessPort,
					Protocol:   corev1.ProtocolTCP,
					Name:       redisAccessPortName,
					TargetPort: intstr.FromInt(redisAccessPort),
				},
			},
			Selector: labels,
		},
	}

}

func generateRedisHeadlessService(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Service {

	name := generateName(svcHeadlessNamePrefix, rc.Name)
	namespace := rc.Namespace

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:     redisHeadlessPort,
					Protocol: corev1.ProtocolTCP,
					Name:     redisHeadlessPortName,
				},
			},
			Selector: labels,
		},
	}
}

func generateRedisConfigMap(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.ConfigMap {
	name := generateName(configMapNamePrefix, rc.Name)
	namespace := rc.Namespace
	conf := `appendonly yes
cluster-enabled yes
cluster-config-file /var/lib/redis/nodes.conf
cluster-node-timeout 5000
dir /var/lib/redis
port 6379`
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Data: map[string]string{
			"redis.conf": conf,
		},
	}
}

func generateName(prefix, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}