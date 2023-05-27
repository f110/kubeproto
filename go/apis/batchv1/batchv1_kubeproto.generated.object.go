package batchv1

import (
	corev1 "go.f110.dev/kubeproto/go/apis/corev1"
	metav1 "go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "apps"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&CronJob{},
		&CronJobList{},
		&Job{},
		&JobList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type CompletionMode string

const (
	CompletionModeNonIndexed CompletionMode = "NonIndexed"
	CompletionModeIndexed    CompletionMode = "Indexed"
)

type ConcurrencyPolicy string

const (
	ConcurrencyPolicyAllow   ConcurrencyPolicy = "Allow"
	ConcurrencyPolicyForbid  ConcurrencyPolicy = "Forbid"
	ConcurrencyPolicyReplace ConcurrencyPolicy = "Replace"
)

type JobConditionType string

const (
	JobConditionTypeSuspended     JobConditionType = "Suspended"
	JobConditionTypeComplete      JobConditionType = "Complete"
	JobConditionTypeFailed        JobConditionType = "Failed"
	JobConditionTypeFailureTarget JobConditionType = "FailureTarget"
)

type PodFailurePolicyAction string

const (
	PodFailurePolicyActionFailJob PodFailurePolicyAction = "FailJob"
	PodFailurePolicyActionIgnore  PodFailurePolicyAction = "Ignore"
	PodFailurePolicyActionCount   PodFailurePolicyAction = "Count"
)

type PodFailurePolicyOnExitCodesOperator string

const (
	PodFailurePolicyOnExitCodesOperatorIn    PodFailurePolicyOnExitCodesOperator = "In"
	PodFailurePolicyOnExitCodesOperatorNotIn PodFailurePolicyOnExitCodesOperator = "NotIn"
)

type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of a cron job, including the schedule.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *CronJobSpec `json:"spec,omitempty"`
	// Current status of a cron job.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status *CronJobStatus `json:"status,omitempty"`
}

func (in *CronJob) DeepCopyInto(out *CronJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(CronJobSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(CronJobStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CronJob) DeepCopy() *CronJob {
	if in == nil {
		return nil
	}
	out := new(CronJob)
	in.DeepCopyInto(out)
	return out
}

func (in *CronJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CronJob `json:"items"`
}

func (in *CronJobList) DeepCopyInto(out *CronJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]CronJob, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CronJobList) DeepCopy() *CronJobList {
	if in == nil {
		return nil
	}
	out := new(CronJobList)
	in.DeepCopyInto(out)
	return out
}

func (in *CronJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of a job.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *JobSpec `json:"spec,omitempty"`
	// Current status of a job.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status *JobStatus `json:"status,omitempty"`
}

func (in *Job) DeepCopyInto(out *Job) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(JobSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(JobStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Job) DeepCopy() *Job {
	if in == nil {
		return nil
	}
	out := new(Job)
	in.DeepCopyInto(out)
	return out
}

func (in *Job) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Job `json:"items"`
}

func (in *JobList) DeepCopyInto(out *JobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Job, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *JobList) DeepCopy() *JobList {
	if in == nil {
		return nil
	}
	out := new(JobList)
	in.DeepCopyInto(out)
	return out
}

func (in *JobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CronJobSpec struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule"`
	// The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// If not specified, this will default to the time zone of the kube-controller-manager process.
	// The set of valid time zone names and the time zone offset is loaded from the system-wide time zone
	// database by the API server during CronJob validation and the controller manager during execution.
	// If no system-wide time zone database can be found a bundled version of the database is used instead.
	// If the time zone name becomes invalid during the lifetime of a CronJob or due to a change in host
	// configuration, the controller will stop creating new new Jobs and will create a system event with the
	// reason UnknownTimeZone.
	// More information can be found in https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#time-zones
	TimeZone string `json:"timeZone,omitempty"`
	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	StartingDeadlineSeconds int64 `json:"startingDeadlineSeconds,omitempty"`
	// Specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`
	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	Suspend bool `json:"suspend,omitempty"`
	// Specifies the job that will be created when executing a CronJob.
	JobTemplate JobTemplateSpec `json:"jobTemplate"`
	// The number of successful finished jobs to retain. Value must be non-negative integer.
	// Defaults to 3.
	SuccessfulJobsHistoryLimit int `json:"successfulJobsHistoryLimit,omitempty"`
	// The number of failed finished jobs to retain. Value must be non-negative integer.
	// Defaults to 1.
	FailedJobsHistoryLimit int `json:"failedJobsHistoryLimit,omitempty"`
}

func (in *CronJobSpec) DeepCopyInto(out *CronJobSpec) {
	*out = *in
	in.JobTemplate.DeepCopyInto(&out.JobTemplate)
}

func (in *CronJobSpec) DeepCopy() *CronJobSpec {
	if in == nil {
		return nil
	}
	out := new(CronJobSpec)
	in.DeepCopyInto(out)
	return out
}

type CronJobStatus struct {
	// A list of pointers to currently running jobs.
	Active []corev1.ObjectReference `json:"active"`
	// Information when was the last time the job was successfully scheduled.
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
	// Information when was the last time the job successfully completed.
	LastSuccessfulTime *metav1.Time `json:"lastSuccessfulTime,omitempty"`
}

func (in *CronJobStatus) DeepCopyInto(out *CronJobStatus) {
	*out = *in
	if in.Active != nil {
		l := make([]corev1.ObjectReference, len(in.Active))
		for i := range in.Active {
			in.Active[i].DeepCopyInto(&l[i])
		}
		out.Active = l
	}
	if in.LastScheduleTime != nil {
		in, out := &in.LastScheduleTime, &out.LastScheduleTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.LastSuccessfulTime != nil {
		in, out := &in.LastSuccessfulTime, &out.LastSuccessfulTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CronJobStatus) DeepCopy() *CronJobStatus {
	if in == nil {
		return nil
	}
	out := new(CronJobStatus)
	in.DeepCopyInto(out)
	return out
}

type JobSpec struct {
	// Specifies the maximum desired number of pods the job should
	// run at any given time. The actual number of pods running in steady state will
	// be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism),
	// i.e. when the work left to do is less than max parallelism.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Parallelism int `json:"parallelism,omitempty"`
	// Specifies the desired number of successfully finished pods the
	// job should be run with.  Setting to null means that the success of any
	// pod signals the success of all pods, and allows parallelism to have any positive
	// value.  Setting to 1 means that parallelism is limited to 1 and the success of that
	// pod signals the success of the job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Completions int `json:"completions,omitempty"`
	// Specifies the duration in seconds relative to the startTime that the job
	// may be continuously active before the system tries to terminate it; value
	// must be positive integer. If a Job is suspended (at creation or through an
	// update), this timer will effectively be stopped and reset when the Job is
	// resumed again.
	ActiveDeadlineSeconds int64 `json:"activeDeadlineSeconds,omitempty"`
	// Specifies the policy of handling failed pods. In particular, it allows to
	// specify the set of actions and conditions which need to be
	// satisfied to take the associated action.
	// If empty, the default behaviour applies - the counter of failed pods,
	// represented by the jobs's .status.failed field, is incremented and it is
	// checked against the backoffLimit. This field cannot be used in combination
	// with restartPolicy=OnFailure.
	// This field is alpha-level. To use this field, you must enable the
	// `JobPodFailurePolicy` feature gate (disabled by default).
	PodFailurePolicy *PodFailurePolicy `json:"podFailurePolicy,omitempty"`
	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	BackoffLimit int `json:"backoffLimit,omitempty"`
	// A label query over pods that should match the pod count.
	// Normally, the system sets this field for you.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	// manualSelector controls generation of pod labels and pod selectors.
	// Leave `manualSelector` unset unless you are certain what you are doing.
	// When false or unset, the system pick labels unique to this job
	// and appends those labels to the pod template.  When true,
	// the user is responsible for picking unique labels and specifying
	// the selector.  Failure to pick a unique label may cause this
	// and other jobs to not function correctly.  However, You may see
	// `manualSelector=true` in jobs that were created with the old `extensions/v1beta1`
	// API.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#specifying-your-own-pod-selector
	ManualSelector bool `json:"manualSelector,omitempty"`
	// Describes the pod that will be created when executing a job.
	// The only allowed template.spec.restartPolicy values are "Never" or "OnFailure".
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Template corev1.PodTemplateSpec `json:"template"`
	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	TTLSecondsAfterFinished int `json:"ttlSecondsAfterFinished,omitempty"`
	// completionMode specifies how Pod completions are tracked. It can be
	// `NonIndexed` (default) or `Indexed`.
	// `NonIndexed` means that the Job is considered complete when there have
	// been .spec.completions successfully completed Pods. Each Pod completion is
	// homologous to each other.
	// `Indexed` means that the Pods of a
	// Job get an associated completion index from 0 to (.spec.completions - 1),
	// available in the annotation batch.kubernetes.io/job-completion-index.
	// The Job is considered complete when there is one successfully completed Pod
	// for each index.
	// When value is `Indexed`, .spec.completions must be specified and
	// `.spec.parallelism` must be less than or equal to 10^5.
	// In addition, The Pod name takes the form
	// `$(job-name)-$(index)-$(random-string)`,
	// the Pod hostname takes the form `$(job-name)-$(index)`.
	// More completion modes can be added in the future.
	// If the Job controller observes a mode that it doesn't recognize, which
	// is possible during upgrades due to version skew, the controller
	// skips updates for the Job.
	CompletionMode CompletionMode `json:"completionMode,omitempty"`
	// suspend specifies whether the Job controller should create Pods or not. If
	// a Job is created with suspend set to true, no Pods are created by the Job
	// controller. If a Job is suspended after creation (i.e. the flag goes from
	// false to true), the Job controller will delete all active Pods associated
	// with this Job. Users must design their workload to gracefully handle this.
	// Suspending a Job will reset the StartTime field of the Job, effectively
	// resetting the ActiveDeadlineSeconds timer too. Defaults to false.
	Suspend bool `json:"suspend,omitempty"`
}

func (in *JobSpec) DeepCopyInto(out *JobSpec) {
	*out = *in
	if in.PodFailurePolicy != nil {
		in, out := &in.PodFailurePolicy, &out.PodFailurePolicy
		*out = new(PodFailurePolicy)
		(*in).DeepCopyInto(*out)
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	in.Template.DeepCopyInto(&out.Template)
}

func (in *JobSpec) DeepCopy() *JobSpec {
	if in == nil {
		return nil
	}
	out := new(JobSpec)
	in.DeepCopyInto(out)
	return out
}

type JobStatus struct {
	// The latest available observations of an object's current state. When a Job
	// fails, one of the conditions will have type "Failed" and status true. When
	// a Job is suspended, one of the conditions will have type "Suspended" and
	// status true; when the Job is resumed, the status of this condition will
	// become false. When a Job is completed, one of the conditions will have
	// type "Complete" and status true.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Conditions []JobCondition `json:"conditions"`
	// Represents time when the job controller started processing a job. When a
	// Job is created in the suspended state, this field is not set until the
	// first time it is resumed. This field is reset every time a Job is resumed
	// from suspension. It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// The completion time is only set when the job finishes successfully.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// The number of pending and running pods.
	Active int `json:"active,omitempty"`
	// The number of pods which reached phase Succeeded.
	Succeeded int `json:"succeeded,omitempty"`
	// The number of pods which reached phase Failed.
	Failed int `json:"failed,omitempty"`
	// completedIndexes holds the completed indexes when .spec.completionMode =
	// "Indexed" in a text format. The indexes are represented as decimal integers
	// separated by commas. The numbers are listed in increasing order. Three or
	// more consecutive numbers are compressed and represented by the first and
	// last element of the series, separated by a hyphen.
	// For example, if the completed indexes are 1, 3, 4, 5 and 7, they are
	// represented as "1,3-5,7".
	CompletedIndexes string `json:"completedIndexes,omitempty"`
	// uncountedTerminatedPods holds the UIDs of Pods that have terminated but
	// the job controller hasn't yet accounted for in the status counters.
	// The job controller creates pods with a finalizer. When a pod terminates
	// (succeeded or failed), the controller does three steps to account for it
	// in the job status:
	// 1. Add the pod UID to the arrays in this field.
	// 2. Remove the pod finalizer.
	// 3. Remove the pod UID from the arrays while increasing the corresponding
	// counter.
	// Old jobs might not be tracked using this field, in which case the field
	// remains null.
	UncountedTerminatedPods *UncountedTerminatedPods `json:"uncountedTerminatedPods,omitempty"`
	// The number of pods which have a Ready condition.
	// This field is beta-level. The job controller populates the field when
	// the feature gate JobReadyPods is enabled (enabled by default).
	Ready int `json:"ready,omitempty"`
}

func (in *JobStatus) DeepCopyInto(out *JobStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]JobCondition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.UncountedTerminatedPods != nil {
		in, out := &in.UncountedTerminatedPods, &out.UncountedTerminatedPods
		*out = new(UncountedTerminatedPods)
		(*in).DeepCopyInto(*out)
	}
}

func (in *JobStatus) DeepCopy() *JobStatus {
	if in == nil {
		return nil
	}
	out := new(JobStatus)
	in.DeepCopyInto(out)
	return out
}

type JobTemplateSpec struct {
	// Standard object's metadata of the jobs created from this template.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	ObjectMeta *metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the job.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *JobSpec `json:"spec,omitempty"`
}

func (in *JobTemplateSpec) DeepCopyInto(out *JobTemplateSpec) {
	*out = *in
	if in.ObjectMeta != nil {
		in, out := &in.ObjectMeta, &out.ObjectMeta
		*out = new(metav1.ObjectMeta)
		(*in).DeepCopyInto(*out)
	}
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(JobSpec)
		(*in).DeepCopyInto(*out)
	}
}

func (in *JobTemplateSpec) DeepCopy() *JobTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(JobTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

type PodFailurePolicy struct {
	// A list of pod failure policy rules. The rules are evaluated in order.
	// Once a rule matches a Pod failure, the remaining of the rules are ignored.
	// When no rule matches the Pod failure, the default handling applies - the
	// counter of pod failures is incremented and it is checked against
	// the backoffLimit. At most 20 elements are allowed.
	Rules []PodFailurePolicyRule `json:"rules"`
}

func (in *PodFailurePolicy) DeepCopyInto(out *PodFailurePolicy) {
	*out = *in
	if in.Rules != nil {
		l := make([]PodFailurePolicyRule, len(in.Rules))
		for i := range in.Rules {
			in.Rules[i].DeepCopyInto(&l[i])
		}
		out.Rules = l
	}
}

func (in *PodFailurePolicy) DeepCopy() *PodFailurePolicy {
	if in == nil {
		return nil
	}
	out := new(PodFailurePolicy)
	in.DeepCopyInto(out)
	return out
}

type JobCondition struct {
	// Type of job condition, Complete or Failed.
	Type JobConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition was checked.
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transit from one status to another.
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// (brief) reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

func (in *JobCondition) DeepCopyInto(out *JobCondition) {
	*out = *in
	if in.LastProbeTime != nil {
		in, out := &in.LastProbeTime, &out.LastProbeTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *JobCondition) DeepCopy() *JobCondition {
	if in == nil {
		return nil
	}
	out := new(JobCondition)
	in.DeepCopyInto(out)
	return out
}

type UncountedTerminatedPods struct {
	// succeeded holds UIDs of succeeded Pods.
	Succeeded []string `json:"succeeded"`
	// failed holds UIDs of failed Pods.
	Failed []string `json:"failed"`
}

func (in *UncountedTerminatedPods) DeepCopyInto(out *UncountedTerminatedPods) {
	*out = *in
	if in.Succeeded != nil {
		t := make([]string, len(in.Succeeded))
		copy(t, in.Succeeded)
		out.Succeeded = t
	}
	if in.Failed != nil {
		t := make([]string, len(in.Failed))
		copy(t, in.Failed)
		out.Failed = t
	}
}

func (in *UncountedTerminatedPods) DeepCopy() *UncountedTerminatedPods {
	if in == nil {
		return nil
	}
	out := new(UncountedTerminatedPods)
	in.DeepCopyInto(out)
	return out
}

type PodFailurePolicyRule struct {
	// Specifies the action taken on a pod failure when the requirements are satisfied.
	// Possible values are:
	// - FailJob: indicates that the pod's job is marked as Failed and all
	// running pods are terminated.
	// - Ignore: indicates that the counter towards the .backoffLimit is not
	// incremented and a replacement pod is created.
	// - Count: indicates that the pod is handled in the default way - the
	// counter towards the .backoffLimit is incremented.
	// Additional values are considered to be added in the future. Clients should
	// react to an unknown action by skipping the rule.
	Action PodFailurePolicyAction `json:"action"`
	// Represents the requirement on the container exit codes.
	OnExitCodes *PodFailurePolicyOnExitCodesRequirement `json:"onExitCodes,omitempty"`
	// Represents the requirement on the pod conditions. The requirement is represented
	// as a list of pod condition patterns. The requirement is satisfied if at
	// least one pattern matches an actual pod condition. At most 20 elements are allowed.
	OnPodConditions []PodFailurePolicyOnPodConditionsPattern `json:"onPodConditions"`
}

func (in *PodFailurePolicyRule) DeepCopyInto(out *PodFailurePolicyRule) {
	*out = *in
	if in.OnExitCodes != nil {
		in, out := &in.OnExitCodes, &out.OnExitCodes
		*out = new(PodFailurePolicyOnExitCodesRequirement)
		(*in).DeepCopyInto(*out)
	}
	if in.OnPodConditions != nil {
		l := make([]PodFailurePolicyOnPodConditionsPattern, len(in.OnPodConditions))
		for i := range in.OnPodConditions {
			in.OnPodConditions[i].DeepCopyInto(&l[i])
		}
		out.OnPodConditions = l
	}
}

func (in *PodFailurePolicyRule) DeepCopy() *PodFailurePolicyRule {
	if in == nil {
		return nil
	}
	out := new(PodFailurePolicyRule)
	in.DeepCopyInto(out)
	return out
}

type PodFailurePolicyOnExitCodesRequirement struct {
	// Restricts the check for exit codes to the container with the
	// specified name. When null, the rule applies to all containers.
	// When specified, it should match one the container or initContainer
	// names in the pod template.
	ContainerName string `json:"containerName,omitempty"`
	// Represents the relationship between the container exit code(s) and the
	// specified values. Containers completed with success (exit code 0) are
	// excluded from the requirement check. Possible values are:
	// - In: the requirement is satisfied if at least one container exit code
	// (might be multiple if there are multiple containers not restricted
	// by the 'containerName' field) is in the set of specified values.
	// - NotIn: the requirement is satisfied if at least one container exit code
	// (might be multiple if there are multiple containers not restricted
	// by the 'containerName' field) is not in the set of specified values.
	// Additional values are considered to be added in the future. Clients should
	// react to an unknown operator by assuming the requirement is not satisfied.
	Operator PodFailurePolicyOnExitCodesOperator `json:"operator"`
	// Specifies the set of values. Each returned container exit code (might be
	// multiple in case of multiple containers) is checked against this set of
	// values with respect to the operator. The list of values must be ordered
	// and must not contain duplicates. Value '0' cannot be used for the In operator.
	// At least one element is required. At most 255 elements are allowed.
	Values []int `json:"values"`
}

func (in *PodFailurePolicyOnExitCodesRequirement) DeepCopyInto(out *PodFailurePolicyOnExitCodesRequirement) {
	*out = *in
	if in.Values != nil {
		t := make([]int, len(in.Values))
		copy(t, in.Values)
		out.Values = t
	}
}

func (in *PodFailurePolicyOnExitCodesRequirement) DeepCopy() *PodFailurePolicyOnExitCodesRequirement {
	if in == nil {
		return nil
	}
	out := new(PodFailurePolicyOnExitCodesRequirement)
	in.DeepCopyInto(out)
	return out
}

type PodFailurePolicyOnPodConditionsPattern struct {
	// Specifies the required Pod condition type. To match a pod condition
	// it is required that specified type equals the pod condition type.
	Type corev1.PodConditionType `json:"type"`
	// Specifies the required Pod condition status. To match a pod condition
	// it is required that the specified status equals the pod condition status.
	// Defaults to True.
	Status corev1.ConditionStatus `json:"status"`
}

func (in *PodFailurePolicyOnPodConditionsPattern) DeepCopyInto(out *PodFailurePolicyOnPodConditionsPattern) {
	*out = *in
}

func (in *PodFailurePolicyOnPodConditionsPattern) DeepCopy() *PodFailurePolicyOnPodConditionsPattern {
	if in == nil {
		return nil
	}
	out := new(PodFailurePolicyOnPodConditionsPattern)
	in.DeepCopyInto(out)
	return out
}
