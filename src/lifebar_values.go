package main

type HealthBarValues struct {
	toplife    float32
	oldlife    float32
	midlife    float32
	midlifeMin float32
	mlifetime  int32
	gethit     bool
}

func newHealthBarValues() *HealthBarValues {
	return &HealthBarValues{oldlife: 1, midlife: 1, midlifeMin: 1}
}

type LifebarValues struct {
	hb [8][]HealthBarValues
}

func (lbv *LifebarValues) clone() (result *LifebarValues) {
	result = &LifebarValues{}
	*result = *lbv

	// Manually copy references that shallow copy poorly, as needed
	// Pointers, slices, maps, functions, channels etc
	for i := range lbv.hb {
		result.hb[i] = make([]HealthBarValues, len(lbv.hb[i]))
		copy(result.hb[i], lbv.hb[i])
	}

	return
}
