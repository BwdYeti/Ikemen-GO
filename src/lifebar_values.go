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

type LifeBarComboValues struct {
	cur, old   int32
	curd, oldd int32
	curp, oldp float32
	resttime   int32
	counterX   float32
	shaketime  int32
	combo      int32
	fakeCombo  int32
}

func newLifeBarComboValues() *LifeBarComboValues {
	return &LifeBarComboValues{}
}

type LifeBarActionValues struct {
	oldleader int
	messages  []*LbMsg
}

func newLifeBarActionValues() *LifeBarActionValues {
	return &LifeBarActionValues{}
}

func (acv *LifeBarActionValues) clone() (result *LifeBarActionValues) {
	result = &LifeBarActionValues{}
	*result = *acv

	// Manually copy references that shallow copy poorly, as needed
	// Pointers, slices, maps, functions, channels etc
	result.messages = make([]*LbMsg, len(acv.messages))
	for i := range acv.messages {
		result.messages[i] = acv.messages[i].clone()
	}

	return
}

type LifebarValues struct {
	hb [8][]HealthBarValues
	co [2]LifeBarComboValues
	ac [2]LifeBarActionValues
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

	for i := range lbv.ac {
		result.ac[i] = *lbv.ac[i].clone()
	}

	return
}
