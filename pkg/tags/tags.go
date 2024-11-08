package tags

const DOInstanceIdPrefix = "hobbyfarm_io_instance_id:"

const DODropletNamePrefix = "digitalocean_hobbyfarm_io_droplet_name:"

// GetTagValue attempts to retrieve a value from a tag formatted in the following manner:
// [tagprefix][tagvalue]
// GetTagValue returns the value of the tag if found, empty string otherwise
func GetTagValue(tags []string, tagPrefix string) string {
	if v, ok := GetTagValues(tags, tagPrefix)[tagPrefix]; !ok {
		return ""
	} else {
		return v
	}
}

// GetTagValues attempts to retrieve values for a list of tag prefixes
// Tags must be defined in the following manner: "{tagprefix}{tagvalue}"
// GetTagVaulues returns a map of tag prefixes to values
func GetTagValues(tags []string, tagPrefixes ...string) (out map[string]string) {
	out = make(map[string]string)
	for _, t := range tags {
		for _, tp := range tagPrefixes {
			if len(t) >= len(tp) {
				if t[:len(tp)] == tp {
					out[tp] = t[len(tp):]
				}
			}
		}
	}

	return
}
