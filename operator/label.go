package operator

// MergeLabels merge all the label maps into one label map
func MergeLabels(allLabels ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, labels := range allLabels {
		if labels != nil {
			for k, v := range labels {
				res[k] = v
			}
		}
	}
	return res
}
