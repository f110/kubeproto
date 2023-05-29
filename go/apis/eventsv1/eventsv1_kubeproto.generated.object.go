package eventsv1

import (
	corev1 "go.f110.dev/kubeproto/go/apis/corev1"
	metav1 "go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "events.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Event{},
		&EventList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type Event struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// eventTime is the time when this Event was first observed. It is required.
	EventTime metav1.MicroTime `json:"eventTime"`
	// series is data about the Event series this event represents or nil if it's a singleton Event.
	Series *EventSeries `json:"series,omitempty"`
	// reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	// This field cannot be empty for new Events.
	ReportingController string `json:"reportingController,omitempty"`
	// reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	ReportingInstance string `json:"reportingInstance,omitempty"`
	// action is what action was taken/failed regarding to the regarding object. It is machine-readable.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	Action string `json:"action,omitempty"`
	// reason is why the action was taken. It is human-readable.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	Reason string `json:"reason,omitempty"`
	// regarding contains the object this Event is about. In most cases it's an Object reporting controller
	// implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because
	// it acts on some changes in a ReplicaSet object.
	Regarding *corev1.ObjectReference `json:"regarding,omitempty"`
	// related is the optional secondary object for more complex actions. E.g. when regarding object triggers
	// a creation or deletion of related object.
	Related *corev1.ObjectReference `json:"related,omitempty"`
	// note is a human-readable description of the status of this operation.
	// Maximal length of the note is 1kB, but libraries should be prepared to
	// handle values up to 64kB.
	Note string `json:"note,omitempty"`
	// type is the type of this event (Normal, Warning), new types could be added in the future.
	// It is machine-readable.
	// This field cannot be empty for new Events.
	Type string `json:"type,omitempty"`
	// deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
	DeprecatedSource *corev1.EventSource `json:"deprecatedSource,omitempty"`
	// deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
	DeprecatedFirstTimestamp *metav1.Time `json:"deprecatedFirstTimestamp,omitempty"`
	// deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
	DeprecatedLastTimestamp *metav1.Time `json:"deprecatedLastTimestamp,omitempty"`
	// deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
	DeprecatedCount int `json:"deprecatedCount,omitempty"`
}

func (in *Event) DeepCopyInto(out *Event) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.EventTime.DeepCopyInto(&out.EventTime)
	if in.Series != nil {
		in, out := &in.Series, &out.Series
		*out = new(EventSeries)
		(*in).DeepCopyInto(*out)
	}
	if in.Regarding != nil {
		in, out := &in.Regarding, &out.Regarding
		*out = new(corev1.ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.Related != nil {
		in, out := &in.Related, &out.Related
		*out = new(corev1.ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.DeprecatedSource != nil {
		in, out := &in.DeprecatedSource, &out.DeprecatedSource
		*out = new(corev1.EventSource)
		(*in).DeepCopyInto(*out)
	}
	if in.DeprecatedFirstTimestamp != nil {
		in, out := &in.DeprecatedFirstTimestamp, &out.DeprecatedFirstTimestamp
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.DeprecatedLastTimestamp != nil {
		in, out := &in.DeprecatedLastTimestamp, &out.DeprecatedLastTimestamp
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Event) DeepCopy() *Event {
	if in == nil {
		return nil
	}
	out := new(Event)
	in.DeepCopyInto(out)
	return out
}

func (in *Event) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type EventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Event `json:"items"`
}

func (in *EventList) DeepCopyInto(out *EventList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Event, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *EventList) DeepCopy() *EventList {
	if in == nil {
		return nil
	}
	out := new(EventList)
	in.DeepCopyInto(out)
	return out
}

func (in *EventList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type EventSeries struct {
	// count is the number of occurrences in this series up to the last heartbeat time.
	Count int `json:"count"`
	// lastObservedTime is the time when last Event from the series was seen before last heartbeat.
	LastObservedTime metav1.MicroTime `json:"lastObservedTime"`
}

func (in *EventSeries) DeepCopyInto(out *EventSeries) {
	*out = *in
	in.LastObservedTime.DeepCopyInto(&out.LastObservedTime)
}

func (in *EventSeries) DeepCopy() *EventSeries {
	if in == nil {
		return nil
	}
	out := new(EventSeries)
	in.DeepCopyInto(out)
	return out
}
