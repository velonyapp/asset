package vo

import "sort"

type Tags struct {
	value map[Tag]struct{}
}

func NewTags(value []string) (Tags, error) {
	tags := make(map[Tag]struct{}, len(value))

	for _, item := range value {
		tag, err := NewTag(item)
		if err != nil {
			return Tags{}, err
		}

		tags[tag] = struct{}{}
	}

	return Tags{
		value: tags,
	}, nil
}

func (t Tags) Values() []Tag {
	value := make([]Tag, 0, len(t.value))

	for tag := range t.value {
		value = append(value, tag)
	}

	sort.Slice(value, func(i, j int) bool {
		return value[i].String() < value[j].String()
	})

	return value
}

func (t Tags) Strings() []string {
	value := make([]string, 0, len(t.value))

	for tag := range t.value {
		value = append(value, tag.String())
	}

	sort.Strings(value)

	return value
}

func (t Tags) Equal(other Tags) bool {
	if len(t.value) != len(other.value) {
		return false
	}

	for tag := range t.value {
		if _, ok := other.value[tag]; !ok {
			return false
		}
	}

	return true
}
