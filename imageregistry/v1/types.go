package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigList is a slice of Config objects.
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Config `json:"items"`
}

const (
	// StorageManagementStateManaged indicates the operator is managing the underlying storage.
	StorageManagementStateManaged = "Managed"
	// StorageManagementStateUnmanaged indicates the operator is not managing the underlying
	// storage.
	StorageManagementStateUnmanaged = "Unmanaged"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the configuration object for a registry instance managed by
// the registry operator
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec ImageRegistrySpec `json:"spec"`
	// +optional
	Status ImageRegistryStatus `json:"status,omitempty"`
}

// ImageRegistrySpec defines the specs for the running registry.
type ImageRegistrySpec struct {
	// operatorSpec allows operator specific configuration to be made.
	operatorv1.OperatorSpec `json:",inline"`
	// managementState indicates whether the registry instance represented
	// by this config instance is under operator management or not.  Valid
	// values are Managed, Unmanaged, and Removed.
	ManagementState operatorv1.ManagementState `json:"managementState"`
	// httpSecret is the value needed by the registry to secure uploads, generated by default.
	// +optional
	HTTPSecret string `json:"httpSecret,omitempty"`
	// proxy defines the proxy to be used when calling master api, upstream
	// registries, etc.
	// +optional
	Proxy ImageRegistryConfigProxy `json:"proxy,omitempty"`
	// storage details for configuring registry storage, e.g. S3 bucket
	// coordinates.
	// +optional
	Storage ImageRegistryConfigStorage `json:"storage,omitempty"`
	// readOnly indicates whether the registry instance should reject attempts
	// to push new images or delete existing ones.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// disableRedirect controls whether to route all data through the Registry,
	// rather than redirecting to the backend.
	// +optional
	DisableRedirect bool `json:"disableRedirect,omitempty"`
	// requests controls how many parallel requests a given registry instance
	// will handle before queuing additional requests.
	// +optional
	Requests ImageRegistryConfigRequests `json:"requests,omitempty"`
	// defaultRoute indicates whether an external facing route for the registry
	// should be created using the default generated hostname.
	// +optional
	DefaultRoute bool `json:"defaultRoute,omitempty"`
	// routes defines additional external facing routes which should be
	// created for the registry.
	// +optional
	Routes []ImageRegistryConfigRoute `json:"routes,omitempty"`
	// replicas determines the number of registry instances to run.
	Replicas int32 `json:"replicas"`
	// logging is deprecated, use logLevel instead.
	// +optional
	Logging int64 `json:"logging,omitempty"`
	// resources defines the resource requests+limits for the registry pod.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// nodeSelector defines the node selection constraints for the registry
	// pod.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// tolerations defines the tolerations for the registry pod.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// rolloutStrategy defines rollout strategy for the image registry
	// deployment.
	// +optional
	// +kubebuilder:validation:Pattern=`^(RollingUpdate|Recreate)$`
	RolloutStrategy string `json:"rolloutStrategy,omitempty"`
	// affinity is a group of node affinity scheduling rules for the image registry pod(s).
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
}

// ImageRegistryStatus reports image registry operational status.
type ImageRegistryStatus struct {
	operatorv1.OperatorStatus `json:",inline"`

	// storageManaged is deprecated, please refer to Storage.managementState
	StorageManaged bool `json:"storageManaged"`
	// storage indicates the current applied storage configuration of the
	// registry.
	Storage ImageRegistryConfigStorage `json:"storage"`
}

// ImageRegistryConfigProxy defines proxy configuration to be used by registry.
type ImageRegistryConfigProxy struct {
	// http defines the proxy to be used by the image registry when
	// accessing HTTP endpoints.
	// +optional
	HTTP string `json:"http,omitempty"`
	// https defines the proxy to be used by the image registry when
	// accessing HTTPS endpoints.
	// +optional
	HTTPS string `json:"https,omitempty"`
	// noProxy defines a comma-separated list of host names that shouldn't
	// go through any proxy.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

// ImageRegistryConfigStorageS3CloudFront holds the configuration
// to use Amazon Cloudfront as the storage middleware in a registry.
// https://docs.docker.com/registry/configuration/#cloudfront
type ImageRegistryConfigStorageS3CloudFront struct {
	// baseURL contains the SCHEME://HOST[/PATH] at which Cloudfront is served.
	BaseURL string `json:"baseURL"`
	// privateKey points to secret containing the private key, provided by AWS.
	PrivateKey corev1.SecretKeySelector `json:"privateKey"`
	// keypairID is key pair ID provided by AWS.
	KeypairID string `json:"keypairID"`
	// duration is the duration of the Cloudfront session.
	// +optional
	// +kubebuilder:validation:Format=duration
	Duration metav1.Duration `json:"duration,omitempty"`
}

// ImageRegistryConfigStorageEmptyDir is an place holder to be used when
// when registry is leveraging ephemeral storage.
type ImageRegistryConfigStorageEmptyDir struct {
}

// ImageRegistryConfigStorageS3 holds the information to configure
// the registry to use the AWS S3 service for backend storage
// https://docs.docker.com/registry/storage-drivers/s3/
type ImageRegistryConfigStorageS3 struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// region is the AWS region in which your bucket exists.
	// Optional, will be set based on the installed AWS Region.
	// +optional
	Region string `json:"region,omitempty"`
	// regionEndpoint is the endpoint for S3 compatible storage services.
	// Optional, defaults based on the Region that is provided.
	// +optional
	RegionEndpoint string `json:"regionEndpoint,omitempty"`
	// encrypt specifies whether the registry stores the image in encrypted
	// format or not.
	// Optional, defaults to false.
	// +optional
	Encrypt bool `json:"encrypt,omitempty"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, Encrypt must be true, or this parameter is ignored.
	// +optional
	KeyID string `json:"keyID,omitempty"`
	// cloudFront configures Amazon Cloudfront as the storage middleware in a
	// registry.
	// +optional
	CloudFront *ImageRegistryConfigStorageS3CloudFront `json:"cloudFront,omitempty"`
	// virtualHostedStyle enables using S3 virtual hosted style bucket paths with
	// a custom RegionEndpoint
	// Optional, defaults to false.
	// +optional
	VirtualHostedStyle bool `json:"virtualHostedStyle"`
}

// ImageRegistryConfigStorageGCS holds GCS configuration.
type ImageRegistryConfigStorageGCS struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// region is the GCS location in which your bucket exists.
	// Optional, will be set based on the installed GCS Region.
	// +optional
	Region string `json:"region,omitempty"`
	// projectID is the Project ID of the GCP project that this bucket should
	// be associated with.
	// +optional
	ProjectID string `json:"projectID,omitempty"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, buckets are encrypted by default on GCP.
	// This allows for the use of a custom encryption key.
	// +optional
	KeyID string `json:"keyID,omitempty"`
}

// ImageRegistryConfigStorageSwift holds the information to configure
// the registry to use the OpenStack Swift service for backend storage
// https://docs.docker.com/registry/storage-drivers/swift/
type ImageRegistryConfigStorageSwift struct {
	// authURL defines the URL for obtaining an authentication token.
	// +optional
	AuthURL string `json:"authURL,omitempty"`
	// authVersion specifies the OpenStack Auth's version.
	// +optional
	AuthVersion string `json:"authVersion,omitempty"`
	// container defines the name of Swift container where to store the
	// registry's data.
	// +optional
	Container string `json:"container,omitempty"`
	// domain specifies Openstack's domain name for Identity v3 API.
	// +optional
	Domain string `json:"domain,omitempty"`
	// domainID specifies Openstack's domain id for Identity v3 API.
	// +optional
	DomainID string `json:"domainID,omitempty"`
	// tenant defines Openstack tenant name to be used by registry.
	// +optional
	Tenant string `json:"tenant,omitempty"`
	// tenant defines Openstack tenant id to be used by registry.
	// +optional
	TenantID string `json:"tenantID,omitempty"`
	// regionName defines Openstack's region in which container exists.
	// +optional
	RegionName string `json:"regionName,omitempty"`
}

// ImageRegistryConfigStoragePVC holds Persistent Volume Claims data to
// be used by the registry.
type ImageRegistryConfigStoragePVC struct {
	// claim defines the Persisent Volume Claim's name to be used.
	// +optional
	Claim string `json:"claim,omitempty"`
}

// ImageRegistryConfigStorageAzure holds the information to configure
// the registry to use Azure Blob Storage for backend storage.
type ImageRegistryConfigStorageAzure struct {
	// accountName defines the account to be used by the registry.
	// +optional
	AccountName string `json:"accountName,omitempty"`
	// container defines Azure's container to be used by registry.
	// +optional
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Pattern=`^[0-9a-z]+(-[0-9a-z]+)*$`
	Container string `json:"container,omitempty"`
	// cloudName is the name of the Azure cloud environment to be used by the
	// registry. If empty, the operator will set it based on the infrastructure
	// object.
	// +optional
	CloudName string `json:"cloudName,omitempty"`
}

// ImageRegistryConfigStorageIBMCOS holds the information to configure
// the registry to use IBM Cloud Object Storage for backend storage.
type ImageRegistryConfigStorageIBMCOS struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// location is the IBM Cloud location in which your bucket exists.
	// Optional, will be set based on the installed IBM Cloud location.
	// +optional
	Location string `json:"location,omitempty"`
	// resourceGroupName is the name of the IBM Cloud resource group that this
	// bucket and its service instance is associated with.
	// Optional, will be set based on the installed IBM Cloud resource group.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// resourceKeyCrn is the CRN of the IBM Cloud resource key that is created
	// for the service instance. Commonly referred as a service credential and
	// must contain HMAC type credentials.
	// Optional, will be computed if not provided.
	// +optional
	ResourceKeyCRN string `json:"resourceKeyCrn,omitempty"`
	// serviceInstanceCrn is the CRN of the IBM Cloud Object Storage service
	// instance that this bucket is associated with.
	// Optional, will be computed if not provided.
	// +optional
	ServiceInstanceCRN string `json:"serviceInstanceCrn,omitempty"`
}

// ImageRegistryConfigStorage describes how the storage should be configured
// for the image registry.
type ImageRegistryConfigStorage struct {
	// emptyDir represents ephemeral storage on the pod's host node.
	// WARNING: this storage cannot be used with more than 1 replica and
	// is not suitable for production use. When the pod is removed from a
	// node for any reason, the data in the emptyDir is deleted forever.
	// +optional
	EmptyDir *ImageRegistryConfigStorageEmptyDir `json:"emptyDir,omitempty"`
	// s3 represents configuration that uses Amazon Simple Storage Service.
	// +optional
	S3 *ImageRegistryConfigStorageS3 `json:"s3,omitempty"`
	// gcs represents configuration that uses Google Cloud Storage.
	// +optional
	GCS *ImageRegistryConfigStorageGCS `json:"gcs,omitempty"`
	// swift represents configuration that uses OpenStack Object Storage.
	// +optional
	Swift *ImageRegistryConfigStorageSwift `json:"swift,omitempty"`
	// pvc represents configuration that uses a PersistentVolumeClaim.
	// +optional
	PVC *ImageRegistryConfigStoragePVC `json:"pvc,omitempty"`
	// azure represents configuration that uses Azure Blob Storage.
	// +optional
	Azure *ImageRegistryConfigStorageAzure `json:"azure,omitempty"`
	// ibmcos represents configuration that uses IBM Cloud Object Storage.
	// +optional
	IBMCOS *ImageRegistryConfigStorageIBMCOS `json:"ibmcos,omitempty"`
	// managementState indicates if the operator manages the underlying
	// storage unit. If Managed the operator will remove the storage when
	// this operator gets Removed.
	// +optional
	// +kubebuilder:validation:Pattern=`^(Managed|Unmanaged)$`
	ManagementState string `json:"managementState,omitempty"`
}

// ImageRegistryConfigRequests defines registry limits on requests read and write.
type ImageRegistryConfigRequests struct {
	// read defines limits for image registry's reads.
	// +optional
	Read ImageRegistryConfigRequestsLimits `json:"read,omitempty"`
	// write defines limits for image registry's writes.
	// +optional
	Write ImageRegistryConfigRequestsLimits `json:"write,omitempty"`
}

// ImageRegistryConfigRequestsLimits holds configuration on the max, enqueued
// and waiting registry's API requests.
type ImageRegistryConfigRequestsLimits struct {
	// maxRunning sets the maximum in flight api requests to the registry.
	// +optional
	MaxRunning int `json:"maxRunning,omitempty"`
	// maxInQueue sets the maximum queued api requests to the registry.
	// +optional
	MaxInQueue int `json:"maxInQueue,omitempty"`
	// maxWaitInQueue sets the maximum time a request can wait in the queue
	// before being rejected.
	// +optional
	// +kubebuilder:validation:Format=duration
	MaxWaitInQueue metav1.Duration `json:"maxWaitInQueue,omitempty"`
}

// ImageRegistryConfigRoute holds information on external route access to image
// registry.
type ImageRegistryConfigRoute struct {
	// name of the route to be created.
	Name string `json:"name"`
	// hostname for the route.
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// secretName points to secret containing the certificates to be used
	// by the route.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}
