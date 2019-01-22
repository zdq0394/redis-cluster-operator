package cluster

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateRedisStatefulset(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *appsv1beta2.StatefulSet {
	name := generateName(statefulsetNamePrefix, rc.Name)
	serviceName := generateName(svcHeadlessNamePrefix, rc.Name)
	namespace := rc.Namespace

	spec := rc.Spec
	redisImage := spec.Redis.Image
	replicas := spec.Redis.Replicas

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			ServiceName: serviceName,
			Replicas:    &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []corev1.Container{
						{
							Name:            "redis",
							Image:           redisImage,
							ImagePullPolicy: "Always",
							Command: []string{
								"redis-server",
							},
							Args: []string{
								"/etc/redis/redis.conf",
								"--protected-mode",
								"no",
							},
							Ports:          getContainerPorts(rc),
							VolumeMounts:   getVolumeMounts(rc),
							Resources:      getResources(rc),
							LivenessProbe:  getLivenessProbe(),
							ReadinessProbe: getReadinessProbe(),
						},
					},
					Volumes: getVolumes(rc),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplates(rc, labels, ownerRefs),
		},
	}
	return ss
}

func getLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: graceTime,
		TimeoutSeconds:      5,
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"sh",
					"-c",
					"redis-cli -h $(hostname) ping",
				},
			},
		},
	}
}

func getReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: graceTime,
		TimeoutSeconds:      5,
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"sh",
					"-c",
					"redis-cli -h $(hostname) ping",
				},
			},
		},
	}
}

func getVolumeMounts(rc *redisv1alpha1.RedisCluster) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "redis-conf",
			MountPath: "/etc/redis",
		},
		{
			Name:      "redis-data",
			MountPath: "/var/lib/redis",
		},
	}
}

func getVolumes(rc *redisv1alpha1.RedisCluster) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "redis-conf",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: generateName(configMapNamePrefix, rc.Name),
					},
				},
			},
		},
	}
}

func getContainerPorts(rc *redisv1alpha1.RedisCluster) []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          redisPortName,
			ContainerPort: redisPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          redisClusterPortName,
			ContainerPort: redisClusterPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}

func generateResourceList(cpu string, memory string) corev1.ResourceList {
	resources := corev1.ResourceList{}
	if cpu != "" {
		resources[corev1.ResourceCPU], _ = resource.ParseQuantity(cpu)
	}
	if memory != "" {
		resources[corev1.ResourceMemory], _ = resource.ParseQuantity(memory)
	}
	return resources
}

func getRequests(resources redisv1alpha1.Resources) corev1.ResourceList {
	return generateResourceList(resources.Requests.CPU, resources.Requests.Memory)
}

func getLimits(resources redisv1alpha1.Resources) corev1.ResourceList {
	return generateResourceList(resources.Limits.CPU, resources.Limits.Memory)
}

func getResources(rc *redisv1alpha1.RedisCluster) corev1.ResourceRequirements {
	spec := rc.Spec
	return corev1.ResourceRequirements{
		Requests: getRequests(spec.Redis.Resources),
		Limits:   getLimits(spec.Redis.Resources),
	}
}

func getVolumeClaimTemplates(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) []corev1.PersistentVolumeClaim {
	storageSize := rc.Spec.Redis.Storage.Size
	storageClassName := rc.Spec.Redis.Storage.StorageClassName
	return []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "redis-data",
				Namespace:       rc.Namespace,
				Labels:          labels,
				OwnerReferences: ownerRefs,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(storageSize),
					},
				},
				StorageClassName: &storageClassName,
			},
		},
	}
}
