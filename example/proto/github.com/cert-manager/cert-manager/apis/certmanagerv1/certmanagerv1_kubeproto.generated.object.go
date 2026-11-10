package certmanagerv1

import (
	"go.f110.dev/kubeproto/example/proto/github.com/cert-manager/cert-manager/apis/acmev1"
	metav1_1 "go.f110.dev/kubeproto/example/proto/github.com/cert-manager/cert-manager/apis/metav1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"maps"
)

const GroupName = "cert-manager.io."

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Certificate{},
		&CertificateList{},
		&CertificateRequest{},
		&CertificateRequestList{},
		&ClusterIssuer{},
		&ClusterIssuerList{},
		&Issuer{},
		&IssuerList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type CertificateConditionType string

const (
	CertificateConditionTypeReady   CertificateConditionType = "Ready"
	CertificateConditionTypeIssuing CertificateConditionType = "Issuing"
)

type CertificateOutputFormatType string

const (
	CertificateOutputFormatTypeDer         CertificateOutputFormatType = "DER"
	CertificateOutputFormatTypeCombinedPEM CertificateOutputFormatType = "CombinedPEM"
)

type CertificateRenewalPolicy string

const (
	CertificateRenewalPolicyRenewBefore CertificateRenewalPolicy = "RenewBefore"
	CertificateRenewalPolicyDisabled    CertificateRenewalPolicy = "Disabled"
)

type CertificateRequestConditionType string

const (
	CertificateRequestConditionTypeReady          CertificateRequestConditionType = "Ready"
	CertificateRequestConditionTypeInvalidRequest CertificateRequestConditionType = "InvalidRequest"
	CertificateRequestConditionTypeApproved       CertificateRequestConditionType = "Approved"
	CertificateRequestConditionTypeDenied         CertificateRequestConditionType = "Denied"
)

type IssuerConditionType string

const (
	IssuerConditionTypeReady IssuerConditionType = "Ready"
)

type KeyUsage string

const (
	KeyUsageSigning           KeyUsage = "signing"
	KeyUsageDigitalSignature  KeyUsage = "digital signature"
	KeyUsageContentCommitment KeyUsage = "content commitment"
	KeyUsageKeyEncipherment   KeyUsage = "key encipherment"
	KeyUsageKeyAgreement      KeyUsage = "key agreement"
	KeyUsageDataEncipherment  KeyUsage = "data encipherment"
	KeyUsageCertSign          KeyUsage = "cert sign"
	KeyUsageCrlSign           KeyUsage = "crl sign"
	KeyUsageEncipherOnly      KeyUsage = "encipher only"
	KeyUsageDecipherOnly      KeyUsage = "decipher only"
	KeyUsageAny               KeyUsage = "any"
	KeyUsageServerAuth        KeyUsage = "server auth"
	KeyUsageClientAuth        KeyUsage = "client auth"
	KeyUsageCodeSigning       KeyUsage = "code signing"
	KeyUsageEmailProtection   KeyUsage = "email protection"
	KeyUsageSMime             KeyUsage = "s/mime"
	KeyUsageIpsecEndSystem    KeyUsage = "ipsec end system"
	KeyUsageIpsecTunnel       KeyUsage = "ipsec tunnel"
	KeyUsageIpsecUser         KeyUsage = "ipsec user"
	KeyUsageTimestamping      KeyUsage = "timestamping"
	KeyUsageOcspSigning       KeyUsage = "ocsp signing"
	KeyUsageMicrosoftSgc      KeyUsage = "microsoft sgc"
	KeyUsageNetscapeSgc       KeyUsage = "netscape sgc"
)

type PKCS12Profile string

const (
	PKCS12ProfileLegacyrc2  PKCS12Profile = "LegacyRC2"
	PKCS12ProfileLegacyDES  PKCS12Profile = "LegacyDES"
	PKCS12ProfileModern2023 PKCS12Profile = "Modern2023"
	PKCS12ProfileModern2026 PKCS12Profile = "Modern2026"
)

type PrivateKeyAlgorithm string

const (
	PrivateKeyAlgorithmRsa     PrivateKeyAlgorithm = "RSA"
	PrivateKeyAlgorithmEcdsa   PrivateKeyAlgorithm = "ECDSA"
	PrivateKeyAlgorithmEd25519 PrivateKeyAlgorithm = "Ed25519"
)

type PrivateKeyEncoding string

const (
	PrivateKeyEncodingPkcs1 PrivateKeyEncoding = "PKCS1"
	PrivateKeyEncodingPkcs8 PrivateKeyEncoding = "PKCS8"
)

type SignatureAlgorithm string

const (
	SignatureAlgorithmSha256withrsa   SignatureAlgorithm = "SHA256WithRSA"
	SignatureAlgorithmSha384withrsa   SignatureAlgorithm = "SHA384WithRSA"
	SignatureAlgorithmSha512withrsa   SignatureAlgorithm = "SHA512WithRSA"
	SignatureAlgorithmEcdsawithsha256 SignatureAlgorithm = "ECDSAWithSHA256"
	SignatureAlgorithmEcdsawithsha384 SignatureAlgorithm = "ECDSAWithSHA384"
	SignatureAlgorithmEcdsawithsha512 SignatureAlgorithm = "ECDSAWithSHA512"
	SignatureAlgorithmPureed25519     SignatureAlgorithm = "PureEd25519"
)

type Certificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired state of the Certificate resource.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec CertificateSpec `json:"spec"`
	// Status of the Certificate.
	// This is set and managed automatically.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status CertificateStatus `json:"status"`
}

func (in *Certificate) DeepCopyInto(out *Certificate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Certificate) DeepCopy() *Certificate {
	if in == nil {
		return nil
	}
	out := new(Certificate)
	in.DeepCopyInto(out)
	return out
}

func (in *Certificate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Certificate) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.Certificate"
}

type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Certificate `json:"items,omitempty"`
}

func (in *CertificateList) DeepCopyInto(out *CertificateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Certificate, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CertificateList) DeepCopy() *CertificateList {
	if in == nil {
		return nil
	}
	out := new(CertificateList)
	in.DeepCopyInto(out)
	return out
}

func (in *CertificateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CertificateList) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.CertificateList"
}

type CertificateRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired state of the CertificateRequest resource.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec CertificateRequestSpec `json:"spec"`
	// Status of the CertificateRequest.
	// This is set and managed automatically.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status CertificateRequestStatus `json:"status"`
}

func (in *CertificateRequest) DeepCopyInto(out *CertificateRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *CertificateRequest) DeepCopy() *CertificateRequest {
	if in == nil {
		return nil
	}
	out := new(CertificateRequest)
	in.DeepCopyInto(out)
	return out
}

func (in *CertificateRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CertificateRequest) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.CertificateRequest"
}

type CertificateRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CertificateRequest `json:"items,omitempty"`
}

func (in *CertificateRequestList) DeepCopyInto(out *CertificateRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]CertificateRequest, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CertificateRequestList) DeepCopy() *CertificateRequestList {
	if in == nil {
		return nil
	}
	out := new(CertificateRequestList)
	in.DeepCopyInto(out)
	return out
}

func (in *CertificateRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CertificateRequestList) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.CertificateRequestList"
}

type ClusterIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Desired state of the ClusterIssuer resource.
	Spec IssuerSpec `json:"spec"`
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status IssuerStatus `json:"status"`
}

func (in *ClusterIssuer) DeepCopyInto(out *ClusterIssuer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *ClusterIssuer) DeepCopy() *ClusterIssuer {
	if in == nil {
		return nil
	}
	out := new(ClusterIssuer)
	in.DeepCopyInto(out)
	return out
}

func (in *ClusterIssuer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ClusterIssuer) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.ClusterIssuer"
}

type ClusterIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ClusterIssuer `json:"items,omitempty"`
}

func (in *ClusterIssuerList) DeepCopyInto(out *ClusterIssuerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ClusterIssuer, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ClusterIssuerList) DeepCopy() *ClusterIssuerList {
	if in == nil {
		return nil
	}
	out := new(ClusterIssuerList)
	in.DeepCopyInto(out)
	return out
}

func (in *ClusterIssuerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ClusterIssuerList) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.ClusterIssuerList"
}

type Issuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Desired state of the Issuer resource.
	Spec IssuerSpec `json:"spec"`
	// Status of the Issuer. This is set and managed automatically.
	Status IssuerStatus `json:"status"`
}

func (in *Issuer) DeepCopyInto(out *Issuer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Issuer) DeepCopy() *Issuer {
	if in == nil {
		return nil
	}
	out := new(Issuer)
	in.DeepCopyInto(out)
	return out
}

func (in *Issuer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Issuer) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.Issuer"
}

type IssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Issuer `json:"items,omitempty"`
}

func (in *IssuerList) DeepCopyInto(out *IssuerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Issuer, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *IssuerList) DeepCopy() *IssuerList {
	if in == nil {
		return nil
	}
	out := new(IssuerList)
	in.DeepCopyInto(out)
	return out
}

func (in *IssuerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *IssuerList) OpenAPIModelName() string {
	return "github.com.cert_manager.cert_manager.pkg.apis.certmanagerv1.IssuerList"
}

type CertificateSpec struct {
	// Requested set of X509 certificate subject attributes.
	// More info: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
	// The common name attribute is specified separately in the `commonName` field.
	// Cannot be set if the `literalSubject` field is set.
	Subject *X509Subject `json:"subject,omitempty"`
	// Requested X.509 certificate subject, represented using the LDAP "String
	// Representation of a Distinguished Name" [1].
	// Important: the LDAP string format also specifies the order of the attributes
	// in the subject, this is important when issuing certs for LDAP authentication.
	// Example: `CN=foo,DC=corp,DC=example,DC=com`
	// More info [1]: https://datatracker.ietf.org/doc/html/rfc4514
	// More info: https://github.com/cert-manager/cert-manager/issues/3203
	// More info: https://github.com/cert-manager/cert-manager/issues/4424
	// Cannot be set if the `subject` or `commonName` field is set.
	LiteralSubject string `json:"literalSubject,omitempty"`
	// Requested common name X509 certificate subject attribute.
	// More info: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
	// NOTE: TLS clients will ignore this value when any subject alternative name is
	// set (see https://tools.ietf.org/html/rfc6125#section-6.4.4).
	// Should have a length of 64 characters or fewer to avoid generating invalid CSRs.
	// Cannot be set if the `literalSubject` field is set.
	CommonName string `json:"commonName,omitempty"`
	// Requested 'duration' (i.e. lifetime) of the Certificate. Note that the
	// issuer may choose to ignore the requested duration, just like any other
	// requested attribute.
	// If unset, this defaults to 90 days.
	// Minimum accepted duration is 1 hour.
	// Value must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration.
	Duration *metav1.Duration `json:"duration,omitempty"`
	// How long before the currently issued certificate's expiry cert-manager should
	// renew the certificate. For example, if a certificate is valid for 60 minutes,
	// and `renewBefore=10m`, cert-manager will begin to attempt to renew the certificate
	// 50 minutes after it was issued (i.e. when there are 10 minutes remaining until
	// the certificate is no longer valid).
	// NOTE: The actual lifetime of the issued certificate is used to determine the
	// renewal time. If an issuer returns a certificate with a different lifetime than
	// the one requested, cert-manager will use the lifetime of the issued certificate.
	// If unset, this defaults to 1/3 of the issued certificate's lifetime.
	// Minimum accepted value is 5 minutes.
	// Value must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration.
	// Cannot be set if the `renewBeforePercentage` field is set.
	RenewBefore *metav1.Duration `json:"renewBefore,omitempty"`
	// `renewBeforePercentage` is like `renewBefore`, except it is a relative percentage
	// rather than an absolute duration. For example, if a certificate is valid for 60
	// minutes, and  `renewBeforePercentage=25`, cert-manager will begin to attempt to
	// renew the certificate 45 minutes after it was issued (i.e. when there are 15
	// minutes (25%) remaining until the certificate is no longer valid).
	// NOTE: The actual lifetime of the issued certificate is used to determine the
	// renewal time. If an issuer returns a certificate with a different lifetime than
	// the one requested, cert-manager will use the lifetime of the issued certificate.
	// Value must be an integer in the range (0,100). The minimum effective
	// `renewBefore` derived from the `renewBeforePercentage` and `duration` fields is 5
	// minutes.
	// Cannot be set if the `renewBefore` field is set.
	RenewBeforePercentage int `json:"renewBeforePercentage,omitempty"`
	// `renewal` allows configuration of how your certificate is renewed. If the policy mentioned is
	// `RenewBefore` then the controller respects `renewBefore` and `renewBeforePercentage`.
	Renewal *CertificateRenewal `json:"renewal,omitempty"`
	// Requested DNS subject alternative names.
	DNSNames []string `json:"dnsNames,omitempty"`
	// Requested IP address subject alternative names.
	IPAddresses []string `json:"ipAddresses,omitempty"`
	// Requested URI subject alternative names.
	URIs []string `json:"uris,omitempty"`
	// `otherNames` is an escape hatch for SAN that allows any type. We currently restrict the support to string like otherNames, cf RFC 5280 p 37
	// Any UTF8 String valued otherName can be passed with by setting the keys oid: x.x.x.x and UTF8Value: somevalue for `otherName`.
	// Most commonly this would be UPN set with oid: 1.3.6.1.4.1.311.20.2.3
	// You should ensure that any OID passed is valid for the UTF8String type as we do not explicitly validate this.
	OtherNames []OtherName `json:"otherNames,omitempty"`
	// Requested email subject alternative names.
	EmailAddresses []string `json:"emailAddresses,omitempty"`
	// Name of the Secret resource that will be automatically created and
	// managed by this Certificate resource. It will be populated with a
	// private key and certificate, signed by the denoted issuer. The Secret
	// resource lives in the same namespace as the Certificate resource.
	SecretName string `json:"secretName"`
	// Defines annotations and labels to be copied to the Certificate's Secret.
	// Labels and annotations on the Secret will be changed as they appear on the
	// SecretTemplate when added or removed. SecretTemplate annotations are added
	// in conjunction with, and cannot overwrite, the base set of annotations
	// cert-manager sets on the Certificate's Secret.
	SecretTemplate *CertificateSecretTemplate `json:"secretTemplate,omitempty"`
	// Additional keystore output formats to be stored in the Certificate's Secret.
	Keystores *CertificateKeystores `json:"keystores,omitempty"`
	// Reference to the issuer responsible for issuing the certificate.
	// If the issuer is namespace-scoped, it must be in the same namespace
	// as the Certificate. If the issuer is cluster-scoped, it can be used
	// from any namespace.
	// The `name` field of the reference must always be specified.
	IssuerRef metav1_1.IssuerReference `json:"issuerRef"`
	// Requested basic constraints isCA value.
	// The isCA value is used to set the `isCA` field on the created CertificateRequest
	// resources. Note that the issuer may choose to ignore the requested isCA value, just
	// like any other requested attribute.
	// If true, this will automatically add the `cert sign` usage to the list
	// of requested `usages`.
	IsCA bool `json:"isCA,omitempty"`
	// Requested key usages and extended key usages.
	// These usages are used to set the `usages` field on the created CertificateRequest
	// resources. If `encodeUsagesInRequest` is unset or set to `true`, the usages
	// will additionally be encoded in the `request` field which contains the CSR blob.
	// If unset, defaults to `digital signature` and `key encipherment`.
	Usages []KeyUsage `json:"usages,omitempty"`
	// Private key options. These include the key algorithm and size, the used
	// encoding and the rotation policy.
	PrivateKey *CertificatePrivateKey `json:"privateKey,omitempty"`
	// Signature algorithm to use.
	// Allowed values for RSA keys: SHA256WithRSA, SHA384WithRSA, SHA512WithRSA.
	// Allowed values for ECDSA keys: ECDSAWithSHA256, ECDSAWithSHA384, ECDSAWithSHA512.
	// Allowed values for Ed25519 keys: PureEd25519.
	SignatureAlgorithm SignatureAlgorithm `json:"signatureAlgorithm,omitempty"`
	// Whether the KeyUsage and ExtKeyUsage extensions should be set in the encoded CSR.
	// This option defaults to true, and should only be disabled if the target
	// issuer does not support CSRs with these X509 KeyUsage/ ExtKeyUsage extensions.
	EncodeUsagesInRequest bool `json:"encodeUsagesInRequest,omitempty"`
	// The maximum number of CertificateRequest revisions that are maintained in
	// the Certificate's history. Each revision represents a single `CertificateRequest`
	// created by this Certificate, either when it was created, renewed, or Spec
	// was changed. Revisions will be removed by oldest first if the number of
	// revisions exceeds this number.
	// If set, revisionHistoryLimit must be a value of `1` or greater.
	// Default value is `1`.
	RevisionHistoryLimit int `json:"revisionHistoryLimit,omitempty"`
	// Defines extra output formats of the private key and signed certificate chain
	// to be written to this Certificate's target Secret.
	AdditionalOutputFormats []CertificateAdditionalOutputFormat `json:"additionalOutputFormats,omitempty"`
	// x.509 certificate NameConstraint extension which MUST NOT be used in a non-CA certificate.
	// More Info: https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.10
	// This is an Alpha Feature and is only enabled with the
	// `--feature-gates=NameConstraints=true` option set on both
	// the controller and webhook components.
	NameConstraints *NameConstraints `json:"nameConstraints,omitempty"`
}

func (in *CertificateSpec) DeepCopyInto(out *CertificateSpec) {
	*out = *in
	if in.Subject != nil {
		in, out := &in.Subject, &out.Subject
		*out = new(X509Subject)
		(*in).DeepCopyInto(*out)
	}
	if in.Duration != nil {
		in, out := &in.Duration, &out.Duration
		*out = new(metav1.Duration)
		(*in).DeepCopyInto(*out)
	}
	if in.RenewBefore != nil {
		in, out := &in.RenewBefore, &out.RenewBefore
		*out = new(metav1.Duration)
		(*in).DeepCopyInto(*out)
	}
	if in.Renewal != nil {
		in, out := &in.Renewal, &out.Renewal
		*out = new(CertificateRenewal)
		(*in).DeepCopyInto(*out)
	}
	if in.DNSNames != nil {
		t := make([]string, len(in.DNSNames))
		copy(t, in.DNSNames)
		out.DNSNames = t
	}
	if in.IPAddresses != nil {
		t := make([]string, len(in.IPAddresses))
		copy(t, in.IPAddresses)
		out.IPAddresses = t
	}
	if in.URIs != nil {
		t := make([]string, len(in.URIs))
		copy(t, in.URIs)
		out.URIs = t
	}
	if in.OtherNames != nil {
		l := make([]OtherName, len(in.OtherNames))
		for i := range in.OtherNames {
			in.OtherNames[i].DeepCopyInto(&l[i])
		}
		out.OtherNames = l
	}
	if in.EmailAddresses != nil {
		t := make([]string, len(in.EmailAddresses))
		copy(t, in.EmailAddresses)
		out.EmailAddresses = t
	}
	if in.SecretTemplate != nil {
		in, out := &in.SecretTemplate, &out.SecretTemplate
		*out = new(CertificateSecretTemplate)
		(*in).DeepCopyInto(*out)
	}
	if in.Keystores != nil {
		in, out := &in.Keystores, &out.Keystores
		*out = new(CertificateKeystores)
		(*in).DeepCopyInto(*out)
	}
	in.IssuerRef.DeepCopyInto(&out.IssuerRef)
	if in.Usages != nil {
		t := make([]KeyUsage, len(in.Usages))
		copy(t, in.Usages)
		out.Usages = t
	}
	if in.PrivateKey != nil {
		in, out := &in.PrivateKey, &out.PrivateKey
		*out = new(CertificatePrivateKey)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalOutputFormats != nil {
		l := make([]CertificateAdditionalOutputFormat, len(in.AdditionalOutputFormats))
		for i := range in.AdditionalOutputFormats {
			in.AdditionalOutputFormats[i].DeepCopyInto(&l[i])
		}
		out.AdditionalOutputFormats = l
	}
	if in.NameConstraints != nil {
		in, out := &in.NameConstraints, &out.NameConstraints
		*out = new(NameConstraints)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateSpec) DeepCopy() *CertificateSpec {
	if in == nil {
		return nil
	}
	out := new(CertificateSpec)
	in.DeepCopyInto(out)
	return out
}

type CertificateStatus struct {
	// List of status conditions to indicate the status of certificates.
	// Known condition types are `Ready` and `Issuing`.
	Conditions []CertificateCondition `json:"conditions,omitempty"`
	// LastFailureTime is set only if the latest issuance for this
	// Certificate failed and contains the time of the failure. If an
	// issuance has failed, the delay till the next issuance will be
	// calculated using formula time.Hour * 2 ^ (failedIssuanceAttempts -
	// 1). If the latest issuance has succeeded this field will be unset.
	LastFailureTime *metav1.Time `json:"lastFailureTime,omitempty"`
	// The time after which the certificate stored in the secret named
	// by this resource in `spec.secretName` is valid.
	NotBefore *metav1.Time `json:"notBefore,omitempty"`
	// The expiration time of the certificate stored in the secret named
	// by this resource in `spec.secretName`.
	NotAfter *metav1.Time `json:"notAfter,omitempty"`
	// RenewalTime is the time at which the certificate will be next
	// renewed.
	// If not set, no upcoming renewal is scheduled.
	RenewalTime *metav1.Time `json:"renewalTime,omitempty"`
	// The current 'revision' of the certificate as issued.
	// When a CertificateRequest resource is created, it will have the
	// `cert-manager.io/certificate-revision` set to one greater than the
	// current value of this field.
	// Upon issuance, this field will be set to the value of the annotation
	// on the CertificateRequest resource used to issue the certificate.
	// Persisting the value on the CertificateRequest resource allows the
	// certificates controller to know whether a request is part of an old
	// issuance or if it is part of the ongoing revision's issuance by
	// checking if the revision value in the annotation is greater than this
	// field.
	Revision int `json:"revision,omitempty"`
	// The name of the Secret resource containing the private key to be used
	// for the next certificate iteration.
	// The keymanager controller will automatically set this field if the
	// `Issuing` condition is set to `True`.
	// It will automatically unset this field when the Issuing condition is
	// not set or False.
	NextPrivateKeySecretName string `json:"nextPrivateKeySecretName,omitempty"`
	// The number of continuous failed issuance attempts up till now. This
	// field gets removed (if set) on a successful issuance and gets set to
	// 1 if unset and an issuance has failed. If an issuance has failed, the
	// delay till the next issuance will be calculated using formula
	// time.Hour * 2 ^ (failedIssuanceAttempts - 1).
	FailedIssuanceAttempts int `json:"failedIssuanceAttempts,omitempty"`
	// ACME stores information that is fetched from the ACME CA server.
	ACME *CertificateACMEStatus `json:"acme,omitempty"`
}

func (in *CertificateStatus) DeepCopyInto(out *CertificateStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]CertificateCondition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.LastFailureTime != nil {
		in, out := &in.LastFailureTime, &out.LastFailureTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.NotBefore != nil {
		in, out := &in.NotBefore, &out.NotBefore
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.NotAfter != nil {
		in, out := &in.NotAfter, &out.NotAfter
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.RenewalTime != nil {
		in, out := &in.RenewalTime, &out.RenewalTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.ACME != nil {
		in, out := &in.ACME, &out.ACME
		*out = new(CertificateACMEStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateStatus) DeepCopy() *CertificateStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateStatus)
	in.DeepCopyInto(out)
	return out
}

type CertificateRequestSpec struct {
	// Requested 'duration' (i.e. lifetime) of the Certificate. Note that the
	// issuer may choose to ignore the requested duration, just like any other
	// requested attribute.
	Duration *metav1.Duration `json:"duration,omitempty"`
	// Reference to the issuer responsible for issuing the certificate.
	// If the issuer is namespace-scoped, it must be in the same namespace
	// as the Certificate. If the issuer is cluster-scoped, it can be used
	// from any namespace.
	// The `name` field of the reference must always be specified.
	IssuerRef metav1_1.IssuerReference `json:"issuerRef"`
	// The PEM-encoded X.509 certificate signing request to be submitted to the
	// issuer for signing.
	// If the CSR has a BasicConstraints extension, its isCA attribute must
	// match the `isCA` value of this CertificateRequest.
	// If the CSR has a KeyUsage extension, its key usages must match the
	// key usages in the `usages` field of this CertificateRequest.
	// If the CSR has a ExtKeyUsage extension, its extended key usages
	// must match the extended key usages in the `usages` field of this
	// CertificateRequest.
	Request []byte `json:"request,omitempty"`
	// Requested basic constraints isCA value. Note that the issuer may choose
	// to ignore the requested isCA value, just like any other requested attribute.
	// NOTE: If the CSR in the `Request` field has a BasicConstraints extension,
	// it must have the same isCA value as specified here.
	// If true, this will automatically add the `cert sign` usage to the list
	// of requested `usages`.
	IsCA bool `json:"isCA,omitempty"`
	// Requested key usages and extended key usages.
	// NOTE: If the CSR in the `Request` field has uses the KeyUsage or
	// ExtKeyUsage extension, these extensions must have the same values
	// as specified here without any additional values.
	// If unset, defaults to `digital signature` and `key encipherment`.
	Usages []KeyUsage `json:"usages,omitempty"`
	// Username contains the name of the user that created the CertificateRequest.
	// Populated by the cert-manager webhook on creation and immutable.
	Username string `json:"username,omitempty"`
	// UID contains the uid of the user that created the CertificateRequest.
	// Populated by the cert-manager webhook on creation and immutable.
	UID string `json:"uid,omitempty"`
	// Groups contains group membership of the user that created the CertificateRequest.
	// Populated by the cert-manager webhook on creation and immutable.
	Groups []string `json:"groups,omitempty"`
}

func (in *CertificateRequestSpec) DeepCopyInto(out *CertificateRequestSpec) {
	*out = *in
	if in.Duration != nil {
		in, out := &in.Duration, &out.Duration
		*out = new(metav1.Duration)
		(*in).DeepCopyInto(*out)
	}
	in.IssuerRef.DeepCopyInto(&out.IssuerRef)
	if in.Usages != nil {
		t := make([]KeyUsage, len(in.Usages))
		copy(t, in.Usages)
		out.Usages = t
	}
	if in.Groups != nil {
		t := make([]string, len(in.Groups))
		copy(t, in.Groups)
		out.Groups = t
	}
}

func (in *CertificateRequestSpec) DeepCopy() *CertificateRequestSpec {
	if in == nil {
		return nil
	}
	out := new(CertificateRequestSpec)
	in.DeepCopyInto(out)
	return out
}

type CertificateRequestStatus struct {
	// List of status conditions to indicate the status of a CertificateRequest.
	// Known condition types are `Ready`, `InvalidRequest`, `Approved` and `Denied`.
	Conditions []CertificateRequestCondition `json:"conditions,omitempty"`
	// The PEM encoded X.509 certificate resulting from the certificate
	// signing request.
	// If not set, the CertificateRequest has either not been completed or has
	// failed. More information on failure can be found by checking the
	// `conditions` field.
	Certificate []byte `json:"certificate,omitempty"`
	// The PEM encoded X.509 certificate of the signer, also known as the Ca
	// (Certificate Authority).
	// This is set on a best-effort basis by different issuers.
	// If not set, the CA is assumed to be unknown/not available.
	CA []byte `json:"ca,omitempty"`
	// FailureTime stores the time that this CertificateRequest failed. This is
	// used to influence garbage collection and back-off.
	FailureTime *metav1.Time `json:"failureTime,omitempty"`
}

func (in *CertificateRequestStatus) DeepCopyInto(out *CertificateRequestStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]CertificateRequestCondition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.FailureTime != nil {
		in, out := &in.FailureTime, &out.FailureTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateRequestStatus) DeepCopy() *CertificateRequestStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateRequestStatus)
	in.DeepCopyInto(out)
	return out
}

type IssuerSpec struct {
	IssuerConfig `json:",inline"`
}

func (in *IssuerSpec) DeepCopyInto(out *IssuerSpec) {
	*out = *in
	out.IssuerConfig = in.IssuerConfig
}

func (in *IssuerSpec) DeepCopy() *IssuerSpec {
	if in == nil {
		return nil
	}
	out := new(IssuerSpec)
	in.DeepCopyInto(out)
	return out
}

type IssuerStatus struct {
	// List of status conditions to indicate the status of a CertificateRequest.
	// Known condition types are `Ready`.
	Conditions []IssuerCondition `json:"conditions,omitempty"`
	// ACME specific status options.
	// This field should only be set if the Issuer is configured to use an ACME
	// server to issue certificates.
	ACME *acmev1.ACMEIssuerStatus `json:"acme,omitempty"`
}

func (in *IssuerStatus) DeepCopyInto(out *IssuerStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]IssuerCondition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.ACME != nil {
		in, out := &in.ACME, &out.ACME
		*out = new(acmev1.ACMEIssuerStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IssuerStatus) DeepCopy() *IssuerStatus {
	if in == nil {
		return nil
	}
	out := new(IssuerStatus)
	in.DeepCopyInto(out)
	return out
}

type X509Subject struct {
	// Organizations to be used on the Certificate.
	Organizations []string `json:"organizations,omitempty"`
	// Countries to be used on the Certificate.
	Countries []string `json:"countries,omitempty"`
	// Organizational Units to be used on the Certificate.
	OrganizationalUnits []string `json:"organizationalUnits,omitempty"`
	// Cities to be used on the Certificate.
	Localities []string `json:"localities,omitempty"`
	// State/Provinces to be used on the Certificate.
	Provinces []string `json:"provinces,omitempty"`
	// Street addresses to be used on the Certificate.
	StreetAddresses []string `json:"streetAddresses,omitempty"`
	// Postal codes to be used on the Certificate.
	PostalCodes []string `json:"postalCodes,omitempty"`
	// Serial number to be used on the Certificate.
	SerialNumber string `json:"serialNumber,omitempty"`
}

func (in *X509Subject) DeepCopyInto(out *X509Subject) {
	*out = *in
	if in.Organizations != nil {
		t := make([]string, len(in.Organizations))
		copy(t, in.Organizations)
		out.Organizations = t
	}
	if in.Countries != nil {
		t := make([]string, len(in.Countries))
		copy(t, in.Countries)
		out.Countries = t
	}
	if in.OrganizationalUnits != nil {
		t := make([]string, len(in.OrganizationalUnits))
		copy(t, in.OrganizationalUnits)
		out.OrganizationalUnits = t
	}
	if in.Localities != nil {
		t := make([]string, len(in.Localities))
		copy(t, in.Localities)
		out.Localities = t
	}
	if in.Provinces != nil {
		t := make([]string, len(in.Provinces))
		copy(t, in.Provinces)
		out.Provinces = t
	}
	if in.StreetAddresses != nil {
		t := make([]string, len(in.StreetAddresses))
		copy(t, in.StreetAddresses)
		out.StreetAddresses = t
	}
	if in.PostalCodes != nil {
		t := make([]string, len(in.PostalCodes))
		copy(t, in.PostalCodes)
		out.PostalCodes = t
	}
}

func (in *X509Subject) DeepCopy() *X509Subject {
	if in == nil {
		return nil
	}
	out := new(X509Subject)
	in.DeepCopyInto(out)
	return out
}

type CertificateRenewal struct {
	// `policy` must be one of `Disabled`, `RenewBefore`.
	Policy CertificateRenewalPolicy `json:"policy,omitempty"`
	// `windows` mentions the behavior of when the renewal must happen.
	Windows []CertificateRenewalWindows `json:"windows,omitempty"`
}

func (in *CertificateRenewal) DeepCopyInto(out *CertificateRenewal) {
	*out = *in
	if in.Windows != nil {
		l := make([]CertificateRenewalWindows, len(in.Windows))
		for i := range in.Windows {
			in.Windows[i].DeepCopyInto(&l[i])
		}
		out.Windows = l
	}
}

func (in *CertificateRenewal) DeepCopy() *CertificateRenewal {
	if in == nil {
		return nil
	}
	out := new(CertificateRenewal)
	in.DeepCopyInto(out)
	return out
}

type OtherName struct {
	// Oid is the object identifier for the otherName SAN.
	// The object identifier must be expressed as a dotted string, for
	// example, "1.2.840.113556.1.4.221".
	OID string `json:"oid,omitempty"`
	// utf8Value is the string value of the otherName SAN.
	// The utf8Value accepts any valid UTF8 string to set as value for the otherName SAN.
	UTF8Value string `json:"utf8Value,omitempty"`
}

func (in *OtherName) DeepCopyInto(out *OtherName) {
	*out = *in
}

func (in *OtherName) DeepCopy() *OtherName {
	if in == nil {
		return nil
	}
	out := new(OtherName)
	in.DeepCopyInto(out)
	return out
}

type CertificateSecretTemplate struct {
	// Annotations is a key value map to be copied to the target Kubernetes Secret.
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a key value map to be copied to the target Kubernetes Secret.
	Labels map[string]string `json:"labels,omitempty"`
}

func (in *CertificateSecretTemplate) DeepCopyInto(out *CertificateSecretTemplate) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
}

func (in *CertificateSecretTemplate) DeepCopy() *CertificateSecretTemplate {
	if in == nil {
		return nil
	}
	out := new(CertificateSecretTemplate)
	in.DeepCopyInto(out)
	return out
}

type CertificateKeystores struct {
	// Jks configures options for storing a JKS keystore in the
	// `spec.secretName` Secret resource.
	JKS *JKSKeystore `json:"jks,omitempty"`
	// Pkcs12 configures options for storing a PKCS12 keystore in the
	// `spec.secretName` Secret resource.
	PKCS12 *PKCS12Keystore `json:"pkcs12,omitempty"`
}

func (in *CertificateKeystores) DeepCopyInto(out *CertificateKeystores) {
	*out = *in
	if in.JKS != nil {
		in, out := &in.JKS, &out.JKS
		*out = new(JKSKeystore)
		(*in).DeepCopyInto(*out)
	}
	if in.PKCS12 != nil {
		in, out := &in.PKCS12, &out.PKCS12
		*out = new(PKCS12Keystore)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateKeystores) DeepCopy() *CertificateKeystores {
	if in == nil {
		return nil
	}
	out := new(CertificateKeystores)
	in.DeepCopyInto(out)
	return out
}

type CertificatePrivateKey struct {
	// RotationPolicy controls how private keys should be regenerated when a
	// re-issuance is being processed.
	// If set to `Never`, a private key will only be generated if one does not
	// already exist in the target `spec.secretName`. If one does exist but it
	// does not have the correct algorithm or size, a warning will be raised
	// to await user intervention.
	// If set to `Always`, a private key matching the specified requirements
	// will be generated whenever a re-issuance occurs.
	// Default is `Always`.
	// The default was changed from `Never` to `Always` in cert-manager >=v1.18.0.
	RotationPolicy string `json:"rotationPolicy"`
	// The private key cryptography standards (PKCS) encoding for this
	// certificate's private key to be encoded in.
	// If provided, allowed values are `PKCS1` and `PKCS8` standing for PKCS#1
	// and PKCS#8, respectively.
	// Defaults to `PKCS1` if not specified.
	Encoding PrivateKeyEncoding `json:"encoding,omitempty"`
	// Algorithm is the private key algorithm of the corresponding private key
	// for this certificate.
	// If provided, allowed values are either `RSA`, `ECDSA` or `Ed25519`.
	// If `algorithm` is specified and `size` is not provided,
	// key size of 2048 will be used for `RSA` key algorithm and
	// key size of 256 will be used for `ECDSA` key algorithm.
	// key size is ignored when using the `Ed25519` key algorithm.
	Algorithm PrivateKeyAlgorithm `json:"algorithm,omitempty"`
	// Size is the key bit size of the corresponding private key for this certificate.
	// If `algorithm` is set to `RSA`, valid values are `2048`, `4096` or `8192`,
	// and will default to `2048` if not specified.
	// If `algorithm` is set to `ECDSA`, valid values are `256`, `384` or `521`,
	// and will default to `256` if not specified.
	// If `algorithm` is set to `Ed25519`, Size is ignored.
	// No other values are allowed.
	Size int `json:"size,omitempty"`
}

func (in *CertificatePrivateKey) DeepCopyInto(out *CertificatePrivateKey) {
	*out = *in
}

func (in *CertificatePrivateKey) DeepCopy() *CertificatePrivateKey {
	if in == nil {
		return nil
	}
	out := new(CertificatePrivateKey)
	in.DeepCopyInto(out)
	return out
}

type CertificateAdditionalOutputFormat struct {
	// Type is the name of the format type that should be written to the
	// Certificate's target Secret.
	Type CertificateOutputFormatType `json:"type"`
}

func (in *CertificateAdditionalOutputFormat) DeepCopyInto(out *CertificateAdditionalOutputFormat) {
	*out = *in
}

func (in *CertificateAdditionalOutputFormat) DeepCopy() *CertificateAdditionalOutputFormat {
	if in == nil {
		return nil
	}
	out := new(CertificateAdditionalOutputFormat)
	in.DeepCopyInto(out)
	return out
}

type NameConstraints struct {
	// if true then the name constraints are marked critical.
	Critical bool `json:"critical,omitempty"`
	// Permitted contains the constraints in which the names must be located.
	Permitted *NameConstraintItem `json:"permitted,omitempty"`
	// Excluded contains the constraints which must be disallowed. Any name matching a
	// restriction in the excluded field is invalid regardless
	// of information appearing in the permitted
	Excluded *NameConstraintItem `json:"excluded,omitempty"`
}

func (in *NameConstraints) DeepCopyInto(out *NameConstraints) {
	*out = *in
	if in.Permitted != nil {
		in, out := &in.Permitted, &out.Permitted
		*out = new(NameConstraintItem)
		(*in).DeepCopyInto(*out)
	}
	if in.Excluded != nil {
		in, out := &in.Excluded, &out.Excluded
		*out = new(NameConstraintItem)
		(*in).DeepCopyInto(*out)
	}
}

func (in *NameConstraints) DeepCopy() *NameConstraints {
	if in == nil {
		return nil
	}
	out := new(NameConstraints)
	in.DeepCopyInto(out)
	return out
}

type CertificateCondition struct {
	// Type of the condition, known values are (`Ready`, `Issuing`).
	Type CertificateConditionType `json:"type"`
	// Status of the condition, one of (`True`, `False`, `Unknown`).
	Status metav1.ConditionStatus `json:"status"`
	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason,omitempty"`
	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message,omitempty"`
	// If set, this represents the .metadata.generation that the condition was
	// set based upon.
	// For instance, if .metadata.generation is currently 12, but the
	// .status.condition[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the Certificate.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func (in *CertificateCondition) DeepCopyInto(out *CertificateCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateCondition) DeepCopy() *CertificateCondition {
	if in == nil {
		return nil
	}
	out := new(CertificateCondition)
	in.DeepCopyInto(out)
	return out
}

type CertificateACMEStatus struct {
	// Ari stores the ACME Renewal Information that is fetched from the ACME server
	// in accordance with RFC 9773. This is only populated if the ARI feature gate is enabled.
	ARI *CertificateACMEARIStatus `json:"ari,omitempty"`
}

func (in *CertificateACMEStatus) DeepCopyInto(out *CertificateACMEStatus) {
	*out = *in
	if in.ARI != nil {
		in, out := &in.ARI, &out.ARI
		*out = new(CertificateACMEARIStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateACMEStatus) DeepCopy() *CertificateACMEStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateACMEStatus)
	in.DeepCopyInto(out)
	return out
}

type CertificateRequestCondition struct {
	// Type of the condition, known values are (`Ready`, `InvalidRequest`,
	// `Approved`, `Denied`).
	Type CertificateRequestConditionType `json:"type"`
	// Status of the condition, one of (`True`, `False`, `Unknown`).
	Status metav1.ConditionStatus `json:"status"`
	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason,omitempty"`
	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message,omitempty"`
}

func (in *CertificateRequestCondition) DeepCopyInto(out *CertificateRequestCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateRequestCondition) DeepCopy() *CertificateRequestCondition {
	if in == nil {
		return nil
	}
	out := new(CertificateRequestCondition)
	in.DeepCopyInto(out)
	return out
}

type IssuerConfig struct {
	// ACME configures this issuer to communicate with a RFC8555 (ACME) server
	// to obtain signed x509 certificates.
	ACME *acmev1.ACMEIssuer `json:"acme,omitempty"`
	// Ca configures this issuer to sign certificates using a signing CA keypair
	// stored in a Secret resource.
	// This is used to build internal PKIs that are managed by cert-manager.
	CA *CAIssuer `json:"ca,omitempty"`
	// Vault configures this issuer to sign certificates using a HashiCorp Vault
	// PKI backend.
	Vault *VaultIssuer `json:"vault,omitempty"`
	// SelfSigned configures this issuer to 'self sign' certificates using the
	// private key used to create the CertificateRequest object.
	SelfSigned *SelfSignedIssuer `json:"selfSigned,omitempty"`
	// Venafi configures this issuer to sign certificates using a CyberArk Certificate Manager Self-Hosted
	// or SaaS policy zone.
	Venafi *VenafiIssuer `json:"venafi,omitempty"`
}

func (in *IssuerConfig) DeepCopyInto(out *IssuerConfig) {
	*out = *in
	if in.ACME != nil {
		in, out := &in.ACME, &out.ACME
		*out = new(acmev1.ACMEIssuer)
		(*in).DeepCopyInto(*out)
	}
	if in.CA != nil {
		in, out := &in.CA, &out.CA
		*out = new(CAIssuer)
		(*in).DeepCopyInto(*out)
	}
	if in.Vault != nil {
		in, out := &in.Vault, &out.Vault
		*out = new(VaultIssuer)
		(*in).DeepCopyInto(*out)
	}
	if in.SelfSigned != nil {
		in, out := &in.SelfSigned, &out.SelfSigned
		*out = new(SelfSignedIssuer)
		(*in).DeepCopyInto(*out)
	}
	if in.Venafi != nil {
		in, out := &in.Venafi, &out.Venafi
		*out = new(VenafiIssuer)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IssuerConfig) DeepCopy() *IssuerConfig {
	if in == nil {
		return nil
	}
	out := new(IssuerConfig)
	in.DeepCopyInto(out)
	return out
}

type IssuerCondition struct {
	// Type of the condition, known values are (`Ready`).
	Type IssuerConditionType `json:"type"`
	// Status of the condition, one of (`True`, `False`, `Unknown`).
	Status metav1.ConditionStatus `json:"status"`
	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason,omitempty"`
	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message,omitempty"`
	// If set, this represents the .metadata.generation that the condition was
	// set based upon.
	// For instance, if .metadata.generation is currently 12, but the
	// .status.condition[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the Issuer.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func (in *IssuerCondition) DeepCopyInto(out *IssuerCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IssuerCondition) DeepCopy() *IssuerCondition {
	if in == nil {
		return nil
	}
	out := new(IssuerCondition)
	in.DeepCopyInto(out)
	return out
}

type CertificateRenewalWindows struct {
	// `timezone` is IANA compliant timezone. For example America/Denver.
	// If this field is not set, timezone is treated as UTC.
	Timezone string `json:"timezone,omitempty"`
	// `windowDuration` is how long the cron definition is active for.
	// Value must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration.
	WindowDuration *metav1.Duration `json:"windowDuration,omitempty"`
	// `cron` is a cron compliant string to allow when the renewal should be allowed. Format is as shown below:
	// * * * * *
	// | | | | |
	// | | | | day of the week (0–6) (Sunday to Saturday;
	// | | | month (1–12)             7 is also Sunday on some systems)
	// | | day of the month (1–31)
	// | hour (0–23)
	// minute (0–59)
	Cron string `json:"cron,omitempty"`
}

func (in *CertificateRenewalWindows) DeepCopyInto(out *CertificateRenewalWindows) {
	*out = *in
	if in.WindowDuration != nil {
		in, out := &in.WindowDuration, &out.WindowDuration
		*out = new(metav1.Duration)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateRenewalWindows) DeepCopy() *CertificateRenewalWindows {
	if in == nil {
		return nil
	}
	out := new(CertificateRenewalWindows)
	in.DeepCopyInto(out)
	return out
}

type JKSKeystore struct {
	// Create enables JKS keystore creation for the Certificate.
	// If true, a file named `keystore.jks` will be created in the target
	// Secret resource, encrypted using the password stored in
	// `passwordSecretRef` or `password`.
	// The keystore file will be updated immediately.
	// If the issuer provided a CA certificate, a file named `truststore.jks`
	// will also be created in the target Secret resource, encrypted using the
	// password stored in `passwordSecretRef`
	// containing the issuing Certificate Authority
	Create bool `json:"create"`
	// Alias specifies the alias of the key in the keystore, required by the JKS format.
	// If not provided, the default alias `certificate` will be used.
	Alias string `json:"alias,omitempty"`
	// PasswordSecretRef is a reference to a non-empty key in a Secret resource
	// containing the password used to encrypt the JKS keystore.
	// Mutually exclusive with password.
	// One of password or passwordSecretRef must provide a password with a non-zero length.
	PasswordSecretRef *metav1_1.SecretKeySelector `json:"passwordSecretRef,omitempty"`
	// Password provides a literal password used to encrypt the JKS keystore.
	// Mutually exclusive with passwordSecretRef.
	// One of password or passwordSecretRef must provide a password with a non-zero length.
	Password string `json:"password,omitempty"`
}

func (in *JKSKeystore) DeepCopyInto(out *JKSKeystore) {
	*out = *in
	if in.PasswordSecretRef != nil {
		in, out := &in.PasswordSecretRef, &out.PasswordSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *JKSKeystore) DeepCopy() *JKSKeystore {
	if in == nil {
		return nil
	}
	out := new(JKSKeystore)
	in.DeepCopyInto(out)
	return out
}

type PKCS12Keystore struct {
	// Create enables PKCS12 keystore creation for the Certificate.
	// If true, a file named `keystore.p12` will be created in the target
	// Secret resource, encrypted using the password stored in
	// `passwordSecretRef` or in `password`.
	// The keystore file will be updated immediately.
	// If the issuer provided a CA certificate, a file named `truststore.p12` will
	// also be created in the target Secret resource, encrypted using the
	// password stored in `passwordSecretRef` containing the issuing Certificate
	// Authority
	Create bool `json:"create"`
	// Profile specifies the key and certificate encryption algorithms and the HMAC algorithm
	// used to create the PKCS12 keystore. Default value is `LegacyRC2` for backward compatibility.
	// If provided, allowed values are:
	// `LegacyRC2`: Deprecated. Not supported by default in OpenSSL 3 or Java 20.
	// `LegacyDES`: Less secure algorithm. Use this option for maximal compatibility.
	// `Modern2023`: Secure algorithm. Use this option in case you have to always use secure algorithms
	// (e.g., because of company policy). Please note that the security of the algorithm is not that important
	// in reality, because the unencrypted certificate and private key are also stored in the Secret.
	// `Modern2026`: Encodes PKCS#12 files using algorithms that are considered modern as of 2026.
	// Private keys and certificates are encrypted using PBES2 with PBKDF2-HMAC-SHA-256 and AES-256-CBC.
	// The MAC algorithm is PBMAC1 with PBKDF2-HMAC-SHA-256 and HMAC-SHA256.
	// Files produced with this profile can be read by OpenSSL 3.4.0 and higher, Java 26 and higher,
	// or with Java using compatible versions of Bouncy Castle. Meets FIPS 140-3 requirements.
	Profile PKCS12Profile `json:"profile,omitempty"`
	// PasswordSecretRef is a reference to a non-empty key in a Secret resource
	// containing the password used to encrypt the PKCS#12 keystore.
	// Mutually exclusive with password.
	// One of password or passwordSecretRef must provide a password with a non-zero length.
	PasswordSecretRef *metav1_1.SecretKeySelector `json:"passwordSecretRef,omitempty"`
	// Password provides a literal password used to encrypt the PKCS#12 keystore.
	// Mutually exclusive with passwordSecretRef.
	// One of password or passwordSecretRef must provide a password with a non-zero length.
	Password string `json:"password,omitempty"`
}

func (in *PKCS12Keystore) DeepCopyInto(out *PKCS12Keystore) {
	*out = *in
	if in.PasswordSecretRef != nil {
		in, out := &in.PasswordSecretRef, &out.PasswordSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *PKCS12Keystore) DeepCopy() *PKCS12Keystore {
	if in == nil {
		return nil
	}
	out := new(PKCS12Keystore)
	in.DeepCopyInto(out)
	return out
}

type NameConstraintItem struct {
	// DNSDomains is a list of DNS domains that are permitted or excluded.
	DNSDomains []string `json:"dnsDomains,omitempty"`
	// IPRanges is a list of IP Ranges that are permitted or excluded.
	// This should be a valid CIDR notation.
	IPRanges []string `json:"ipRanges,omitempty"`
	// EmailAddresses is a list of Email Addresses that are permitted or excluded.
	EmailAddresses []string `json:"emailAddresses,omitempty"`
	// URIDomains is a list of URI domains that are permitted or excluded.
	URIDomains []string `json:"uriDomains,omitempty"`
}

func (in *NameConstraintItem) DeepCopyInto(out *NameConstraintItem) {
	*out = *in
	if in.DNSDomains != nil {
		t := make([]string, len(in.DNSDomains))
		copy(t, in.DNSDomains)
		out.DNSDomains = t
	}
	if in.IPRanges != nil {
		t := make([]string, len(in.IPRanges))
		copy(t, in.IPRanges)
		out.IPRanges = t
	}
	if in.EmailAddresses != nil {
		t := make([]string, len(in.EmailAddresses))
		copy(t, in.EmailAddresses)
		out.EmailAddresses = t
	}
	if in.URIDomains != nil {
		t := make([]string, len(in.URIDomains))
		copy(t, in.URIDomains)
		out.URIDomains = t
	}
}

func (in *NameConstraintItem) DeepCopy() *NameConstraintItem {
	if in == nil {
		return nil
	}
	out := new(NameConstraintItem)
	in.DeepCopyInto(out)
	return out
}

type CertificateACMEARIStatus struct {
	// SuggestedWindow is the suggested renewal window as returned by the ACME server in accordance with RFC 9773.
	SuggestedWindow *ACMERenewalWindow `json:"suggestedWindow,omitempty"`
	// ExplanationURL is a human-readable URL that may explain why the suggested window
	// has its current value.
	ExplanationURL string `json:"explanationURL,omitempty"`
	// LastChecked is the time at which the ACME server was last checked for renewal information.
	LastChecked *metav1.Time `json:"lastChecked,omitempty"`
	// NextCheck is the time at which the ACME server will next be checked for renewal information.
	NextCheck *metav1.Time `json:"nextCheck,omitempty"`
	// LastError is the last error encountered when checking the ACME server for renewal information, if any.
	LastError string `json:"lastError,omitempty"`
	// CertID is the ARI CertID (RFC 9773) of the certificate this renewal information was
	// fetched for. This is used to determine if we need to re-fetch the renewal information
	// as changed cert id means that ARI fetched before is stale.
	CertID string `json:"certID,omitempty"`
}

func (in *CertificateACMEARIStatus) DeepCopyInto(out *CertificateACMEARIStatus) {
	*out = *in
	if in.SuggestedWindow != nil {
		in, out := &in.SuggestedWindow, &out.SuggestedWindow
		*out = new(ACMERenewalWindow)
		(*in).DeepCopyInto(*out)
	}
	if in.LastChecked != nil {
		in, out := &in.LastChecked, &out.LastChecked
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.NextCheck != nil {
		in, out := &in.NextCheck, &out.NextCheck
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CertificateACMEARIStatus) DeepCopy() *CertificateACMEARIStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateACMEARIStatus)
	in.DeepCopyInto(out)
	return out
}

type CAIssuer struct {
	// SecretName is the name of the secret used to sign Certificates issued
	// by this Issuer.
	SecretName string `json:"secretName"`
	// The CRL distribution points is an X.509 v3 certificate extension which identifies
	// the location of the CRL from which the revocation of this certificate can be checked.
	// If not set, certificates will be issued without distribution points set.
	CRLDistributionPoints []string `json:"crlDistributionPoints,omitempty"`
	// The OCSP server list is an X.509 v3 extension that defines a list of
	// URLs of OCSP responders. The OCSP responders can be queried for the
	// revocation status of an issued certificate. If not set, the
	// certificate will be issued with no OCSP servers set. For example, an
	// OCSP server URL could be "http://ocsp.int-x3.letsencrypt.org".
	OCSPServers []string `json:"ocspServers,omitempty"`
	// IssuingCertificateURLs is a list of URLs which this issuer should embed into certificates
	// it creates. See https://www.rfc-editor.org/rfc/rfc5280#section-4.2.2.1 for more details.
	// As an example, such a URL might be "http://ca.domain.com/ca.crt".
	IssuingCertificateURLs []string `json:"issuingCertificateURLs,omitempty"`
}

func (in *CAIssuer) DeepCopyInto(out *CAIssuer) {
	*out = *in
	if in.CRLDistributionPoints != nil {
		t := make([]string, len(in.CRLDistributionPoints))
		copy(t, in.CRLDistributionPoints)
		out.CRLDistributionPoints = t
	}
	if in.OCSPServers != nil {
		t := make([]string, len(in.OCSPServers))
		copy(t, in.OCSPServers)
		out.OCSPServers = t
	}
	if in.IssuingCertificateURLs != nil {
		t := make([]string, len(in.IssuingCertificateURLs))
		copy(t, in.IssuingCertificateURLs)
		out.IssuingCertificateURLs = t
	}
}

func (in *CAIssuer) DeepCopy() *CAIssuer {
	if in == nil {
		return nil
	}
	out := new(CAIssuer)
	in.DeepCopyInto(out)
	return out
}

type VaultIssuer struct {
	// Auth configures how cert-manager authenticates with the Vault server.
	Auth VaultAuth `json:"auth"`
	// Server is the connection address for the Vault server, e.g: "https://vault.example.com:8200".
	Server string `json:"server"`
	// ServerName is used to verify the hostname on the returned certificates
	// by the Vault server.
	ServerName string `json:"serverName,omitempty"`
	// Path is the mount path of the Vault PKI backend's `sign` endpoint, e.g:
	// "my_pki_mount/sign/my-role-name".
	Path string `json:"path"`
	// Name of the vault namespace. Namespaces is a set of features within Vault Enterprise that allows Vault environments to support Secure Multi-tenancy. e.g: "ns1"
	// More about namespaces can be found here https://www.vaultproject.io/docs/enterprise/namespaces
	Namespace string `json:"namespace,omitempty"`
	// Base64-encoded bundle of PEM CAs which will be used to validate the certificate
	// chain presented by Vault. Only used if using HTTPS to connect to Vault and
	// ignored for HTTP connections.
	// Mutually exclusive with CABundleSecretRef.
	// If neither CABundle nor CABundleSecretRef are defined, the certificate bundle in
	// the cert-manager controller container is used to validate the TLS connection.
	CABundle []byte `json:"caBundle,omitempty"`
	// Reference to a Secret containing a bundle of PEM-encoded CAs to use when
	// verifying the certificate chain presented by Vault when using HTTPS.
	// Mutually exclusive with CABundle.
	// If neither CABundle nor CABundleSecretRef are defined, the certificate bundle in
	// the cert-manager controller container is used to validate the TLS connection.
	// If no key for the Secret is specified, cert-manager will default to 'ca.crt'.
	CABundleSecretRef *metav1_1.SecretKeySelector `json:"caBundleSecretRef,omitempty"`
	// Reference to a Secret containing a PEM-encoded Client Certificate to use when the
	// Vault server requires mTLS.
	ClientCertSecretRef *metav1_1.SecretKeySelector `json:"clientCertSecretRef,omitempty"`
	// Reference to a Secret containing a PEM-encoded Client Private Key to use when the
	// Vault server requires mTLS.
	ClientKeySecretRef *metav1_1.SecretKeySelector `json:"clientKeySecretRef,omitempty"`
}

func (in *VaultIssuer) DeepCopyInto(out *VaultIssuer) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
	if in.CABundleSecretRef != nil {
		in, out := &in.CABundleSecretRef, &out.CABundleSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ClientCertSecretRef != nil {
		in, out := &in.ClientCertSecretRef, &out.ClientCertSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ClientKeySecretRef != nil {
		in, out := &in.ClientKeySecretRef, &out.ClientKeySecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VaultIssuer) DeepCopy() *VaultIssuer {
	if in == nil {
		return nil
	}
	out := new(VaultIssuer)
	in.DeepCopyInto(out)
	return out
}

type SelfSignedIssuer struct {
	// The CRL distribution points is an X.509 v3 certificate extension which identifies
	// the location of the CRL from which the revocation of this certificate can be checked.
	// If not set certificate will be issued without CDP. Values are strings.
	CRLDistributionPoints []string `json:"crlDistributionPoints,omitempty"`
}

func (in *SelfSignedIssuer) DeepCopyInto(out *SelfSignedIssuer) {
	*out = *in
	if in.CRLDistributionPoints != nil {
		t := make([]string, len(in.CRLDistributionPoints))
		copy(t, in.CRLDistributionPoints)
		out.CRLDistributionPoints = t
	}
}

func (in *SelfSignedIssuer) DeepCopy() *SelfSignedIssuer {
	if in == nil {
		return nil
	}
	out := new(SelfSignedIssuer)
	in.DeepCopyInto(out)
	return out
}

type VenafiIssuer struct {
	// Zone is the Certificate Manager Policy Zone to use for this issuer.
	// All requests made to the Certificate Manager platform will be restricted by the named
	// zone policy.
	// This field is required.
	Zone string `json:"zone"`
	// Tpp specifies CyberArk Certificate Manager Self-Hosted configuration settings.
	// Only one of CyberArk Certificate Manager may be specified.
	TPP *VenafiTPP `json:"tpp,omitempty"`
	// Cloud specifies the CyberArk Certificate Manager SaaS configuration settings.
	// Only one of CyberArk Certificate Manager may be specified.
	Cloud *VenafiCloud `json:"cloud,omitempty"`
	// Ngts specifies Palo Alto Networks Next Generation Trust Services (NGTS) configuration
	// using OAuth 2.0 Client Credentials. Only one of tpp, cloud, or ngts may be specified.
	NGTS *VenafiNGTS `json:"ngts,omitempty"`
}

func (in *VenafiIssuer) DeepCopyInto(out *VenafiIssuer) {
	*out = *in
	if in.TPP != nil {
		in, out := &in.TPP, &out.TPP
		*out = new(VenafiTPP)
		(*in).DeepCopyInto(*out)
	}
	if in.Cloud != nil {
		in, out := &in.Cloud, &out.Cloud
		*out = new(VenafiCloud)
		(*in).DeepCopyInto(*out)
	}
	if in.NGTS != nil {
		in, out := &in.NGTS, &out.NGTS
		*out = new(VenafiNGTS)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VenafiIssuer) DeepCopy() *VenafiIssuer {
	if in == nil {
		return nil
	}
	out := new(VenafiIssuer)
	in.DeepCopyInto(out)
	return out
}

type ACMERenewalWindow struct {
	// Start is the start of the suggested renewal window.
	Start *metav1.Time `json:"start,omitempty"`
	// End is the end of the suggested renewal window.
	End *metav1.Time `json:"end,omitempty"`
}

func (in *ACMERenewalWindow) DeepCopyInto(out *ACMERenewalWindow) {
	*out = *in
	if in.Start != nil {
		in, out := &in.Start, &out.Start
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.End != nil {
		in, out := &in.End, &out.End
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ACMERenewalWindow) DeepCopy() *ACMERenewalWindow {
	if in == nil {
		return nil
	}
	out := new(ACMERenewalWindow)
	in.DeepCopyInto(out)
	return out
}

type VaultAuth struct {
	// TokenSecretRef authenticates with Vault by presenting a token.
	TokenSecretRef *metav1_1.SecretKeySelector `json:"tokenSecretRef,omitempty"`
	// AppRole authenticates with Vault using the App Role auth mechanism,
	// with the role and secret stored in a Kubernetes Secret resource.
	AppRole *VaultAppRole `json:"appRole,omitempty"`
	// ClientCertificate authenticates with Vault by presenting a client
	// certificate during the request's TLS handshake.
	// Works only when using HTTPS protocol.
	ClientCertificate *VaultClientCertificateAuth `json:"clientCertificate,omitempty"`
	// Kubernetes authenticates with Vault by passing the ServiceAccount
	// token stored in the named Secret resource to the Vault server.
	Kubernetes *VaultKubernetesAuth `json:"kubernetes,omitempty"`
	// AWS authenticates with Vault using AWS IAM authentication.
	// This allows authentication using IAM roles for service accounts (IRSA),
	// EKS Pod Identity (PIA), or ambient credentials (EC2 instance profiles, ECS task role).
	AWS *VaultAWSAuth `json:"aws,omitempty"`
}

func (in *VaultAuth) DeepCopyInto(out *VaultAuth) {
	*out = *in
	if in.TokenSecretRef != nil {
		in, out := &in.TokenSecretRef, &out.TokenSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.AppRole != nil {
		in, out := &in.AppRole, &out.AppRole
		*out = new(VaultAppRole)
		(*in).DeepCopyInto(*out)
	}
	if in.ClientCertificate != nil {
		in, out := &in.ClientCertificate, &out.ClientCertificate
		*out = new(VaultClientCertificateAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.Kubernetes != nil {
		in, out := &in.Kubernetes, &out.Kubernetes
		*out = new(VaultKubernetesAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(VaultAWSAuth)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VaultAuth) DeepCopy() *VaultAuth {
	if in == nil {
		return nil
	}
	out := new(VaultAuth)
	in.DeepCopyInto(out)
	return out
}

type VenafiTPP struct {
	// URL is the base URL for the vedsdk endpoint of the CyberArk Certificate Manager Self-Hosted instance,
	// for example: "https://tpp.example.com/vedsdk".
	URL string `json:"url"`
	// CredentialsRef is a reference to a Secret containing the CyberArk Certificate Manager Self-Hosted API credentials.
	// The secret must contain the key 'access-token' for the Access Token Authentication,
	// or two keys, 'username' and 'password' for the API Keys Authentication.
	CredentialsRef metav1_1.LocalObjectReference `json:"credentialsRef"`
	// Base64-encoded bundle of PEM CAs which will be used to validate the certificate
	// chain presented by the CyberArk Certificate Manager Self-Hosted server. Only used if using HTTPS; ignored for HTTP.
	// If undefined, the certificate bundle in the cert-manager controller container
	// is used to validate the chain.
	CABundle []byte `json:"caBundle,omitempty"`
	// Reference to a Secret containing a base64-encoded bundle of PEM CAs
	// which will be used to validate the certificate chain presented by the CyberArk Certificate Manager Self-Hosted server.
	// Only used if using HTTPS; ignored for HTTP. Mutually exclusive with CABundle.
	// If neither CABundle nor CABundleSecretRef is defined, the certificate bundle in
	// the cert-manager controller container is used to validate the TLS connection.
	CABundleSecretRef *metav1_1.SecretKeySelector `json:"caBundleSecretRef,omitempty"`
}

func (in *VenafiTPP) DeepCopyInto(out *VenafiTPP) {
	*out = *in
	in.CredentialsRef.DeepCopyInto(&out.CredentialsRef)
	if in.CABundleSecretRef != nil {
		in, out := &in.CABundleSecretRef, &out.CABundleSecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VenafiTPP) DeepCopy() *VenafiTPP {
	if in == nil {
		return nil
	}
	out := new(VenafiTPP)
	in.DeepCopyInto(out)
	return out
}

type VenafiCloud struct {
	// URL is the base URL for CyberArk Certificate Manager SaaS.
	// Defaults to "https://api.venafi.cloud/".
	URL string `json:"url,omitempty"`
	// APITokenSecretRef is a secret key selector for the CyberArk Certificate Manager SaaS API token.
	APITokenSecretRef metav1_1.SecretKeySelector `json:"apiTokenSecretRef"`
}

func (in *VenafiCloud) DeepCopyInto(out *VenafiCloud) {
	*out = *in
	in.APITokenSecretRef.DeepCopyInto(&out.APITokenSecretRef)
}

func (in *VenafiCloud) DeepCopy() *VenafiCloud {
	if in == nil {
		return nil
	}
	out := new(VenafiCloud)
	in.DeepCopyInto(out)
	return out
}

type VenafiNGTS struct {
	// URL is the base URL for the NGTS API endpoint.
	// Defaults to "https://api.strata.paloaltonetworks.com/ngts" if not set.
	URL string `json:"url,omitempty"`
	// TokenEndpoint is the OAuth 2.0 token endpoint URL used to obtain access tokens,
	// for example "https://auth.apps.paloaltonetworks.com/oauth2/access_token".
	// Defaults to "https://auth.apps.paloaltonetworks.com/oauth2/access_token" if not set.
	TokenEndpoint string `json:"tokenEndpoint,omitempty"`
	// Tsgid is the Tenant Service Group ID used to scope the OAuth 2.0 access token,
	// for example "1234567890". The tsg_id: prefix is added automatically.
	// This field is required.
	TSGID string `json:"tsgID"`
	// CredentialsRef is a reference to a Kubernetes Secret containing the OAuth 2.0
	// Client ID and Client Secret. The secret must contain the keys 'client-id' and
	// 'client-secret'.
	CredentialsRef metav1_1.LocalObjectReference `json:"credentialsRef"`
}

func (in *VenafiNGTS) DeepCopyInto(out *VenafiNGTS) {
	*out = *in
	in.CredentialsRef.DeepCopyInto(&out.CredentialsRef)
}

func (in *VenafiNGTS) DeepCopy() *VenafiNGTS {
	if in == nil {
		return nil
	}
	out := new(VenafiNGTS)
	in.DeepCopyInto(out)
	return out
}

type VaultAppRole struct {
	// Path where the App Role authentication backend is mounted in Vault, e.g:
	// "approle"
	Path string `json:"path"`
	// RoleID configured in the App Role authentication backend when setting
	// up the authentication backend in Vault.
	RoleId string `json:"roleId"`
	// Reference to a key in a Secret that contains the App Role secret used
	// to authenticate with Vault.
	// The `key` field must be specified and denotes which entry within the Secret
	// resource is used as the app role secret.
	SecretRef metav1_1.SecretKeySelector `json:"secretRef"`
}

func (in *VaultAppRole) DeepCopyInto(out *VaultAppRole) {
	*out = *in
	in.SecretRef.DeepCopyInto(&out.SecretRef)
}

func (in *VaultAppRole) DeepCopy() *VaultAppRole {
	if in == nil {
		return nil
	}
	out := new(VaultAppRole)
	in.DeepCopyInto(out)
	return out
}

type VaultClientCertificateAuth struct {
	// The Vault mountPath here is the mount path to use when authenticating with
	// Vault. For example, setting a value to `/v1/auth/foo`, will use the path
	// `/v1/auth/foo/login` to authenticate with Vault. If unspecified, the
	// default value "/v1/auth/cert" will be used.
	Path string `json:"mountPath,omitempty"`
	// Reference to Kubernetes Secret of type "kubernetes.io/tls" (hence containing
	// tls.crt and tls.key) used to authenticate to Vault using TLS client
	// authentication.
	SecretName string `json:"secretName,omitempty"`
	// Name of the certificate role to authenticate against.
	// If not set, matching any certificate role, if available.
	Name string `json:"name,omitempty"`
}

func (in *VaultClientCertificateAuth) DeepCopyInto(out *VaultClientCertificateAuth) {
	*out = *in
}

func (in *VaultClientCertificateAuth) DeepCopy() *VaultClientCertificateAuth {
	if in == nil {
		return nil
	}
	out := new(VaultClientCertificateAuth)
	in.DeepCopyInto(out)
	return out
}

type VaultKubernetesAuth struct {
	// The Vault mountPath here is the mount path to use when authenticating with
	// Vault. For example, setting a value to `/v1/auth/foo`, will use the path
	// `/v1/auth/foo/login` to authenticate with Vault. If unspecified, the
	// default value "/v1/auth/kubernetes" will be used.
	Path string `json:"mountPath,omitempty"`
	// The required Secret field containing a Kubernetes ServiceAccount JWT used
	// for authenticating with Vault. Use of 'ambient credentials' is not
	// supported.
	SecretRef *metav1_1.SecretKeySelector `json:"secretRef,omitempty"`
	// A reference to a service account that will be used to request a bound
	// token (also known as "projected token"). Compared to using "secretRef",
	// using this field means that you don't rely on statically bound tokens. To
	// use this field, you must configure an RBAC rule to let cert-manager
	// request a token.
	ServiceAccountRef *ServiceAccountRef `json:"serviceAccountRef,omitempty"`
	// A required field containing the Vault Role to assume. A Role binds a
	// Kubernetes ServiceAccount with a set of Vault policies.
	Role string `json:"role"`
}

func (in *VaultKubernetesAuth) DeepCopyInto(out *VaultKubernetesAuth) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(metav1_1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ServiceAccountRef != nil {
		in, out := &in.ServiceAccountRef, &out.ServiceAccountRef
		*out = new(ServiceAccountRef)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VaultKubernetesAuth) DeepCopy() *VaultKubernetesAuth {
	if in == nil {
		return nil
	}
	out := new(VaultKubernetesAuth)
	in.DeepCopyInto(out)
	return out
}

type VaultAWSAuth struct {
	// The Vault mountPath here is the mount path to use when authenticating with
	// Vault. For example, setting a value to `/v1/auth/foo`, will use the path
	// `/v1/auth/foo/login` to authenticate with Vault. If unspecified, the
	// default value "/v1/auth/aws" will be used.
	MountPath string `json:"mountPath,omitempty"`
	// A required field containing the Vault Role to assume when authenticating.
	Role string `json:"role"`
	// The AWS region to use for authentication. If not specified, the region
	// will be determined from AWS_REGION or AWS_DEFAULT_REGION environment
	// variables, falling back to "us-east-1" if not set.
	Region string `json:"region,omitempty"`
	// A reference to a service account that will be used to request a web identity
	// token for IRSA (IAM Roles for Service Accounts) authentication.
	ServiceAccountRef *ServiceAccountRef `json:"serviceAccountRef,omitempty"`
	// The ARN of the AWS IAM role to assume using the Kubernetes service account
	// token. Required when using IRSA (serviceAccountRef is set).
	// This role must have a trust policy that allows the OIDC provider to assume it.
	IAMRoleARN string `json:"iamRoleArn,omitempty"`
	// The Vault header value to include in the STS signing request.
	// This is used to prevent replay attacks.
	VaultHeaderValue string `json:"vaultHeaderValue,omitempty"`
}

func (in *VaultAWSAuth) DeepCopyInto(out *VaultAWSAuth) {
	*out = *in
	if in.ServiceAccountRef != nil {
		in, out := &in.ServiceAccountRef, &out.ServiceAccountRef
		*out = new(ServiceAccountRef)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VaultAWSAuth) DeepCopy() *VaultAWSAuth {
	if in == nil {
		return nil
	}
	out := new(VaultAWSAuth)
	in.DeepCopyInto(out)
	return out
}

type ServiceAccountRef struct {
	// Name of the ServiceAccount used to request a token.
	Name string `json:"name"`
	// TokenAudiences is an optional list of extra audiences to include in the token passed to Vault.
	// The default audiences are always included in the token.
	TokenAudiences []string `json:"audiences,omitempty"`
}

func (in *ServiceAccountRef) DeepCopyInto(out *ServiceAccountRef) {
	*out = *in
	if in.TokenAudiences != nil {
		t := make([]string, len(in.TokenAudiences))
		copy(t, in.TokenAudiences)
		out.TokenAudiences = t
	}
}

func (in *ServiceAccountRef) DeepCopy() *ServiceAccountRef {
	if in == nil {
		return nil
	}
	out := new(ServiceAccountRef)
	in.DeepCopyInto(out)
	return out
}
