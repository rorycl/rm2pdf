/*
LocalPenColour set pen colours from an rgba, rgb or hex colour string or from
an image/color colourname as set out in image/colornames.
*/

package penconfig

import "testing"

// TestLocalPenColour
func TestLocalPenColour(t *testing.T) {

	for i, test := range []struct {
		title string
		cName string
		isErr bool
	}{
		{
			title: "blue",
			cName: "blue",
			isErr: false,
		},
		{
			title: "nonsense name",
			cName: "nonsense name",
			isErr: true,
		},
		{
			title: "bluish",
			cName: "#2136D8",
			isErr: false,
		},
		{
			title: "badhash",
			cName: "#2136D8XX",
			isErr: true,
		},
		{
			title: "purplish",
			cName: "rgb(172,33,216)",
			isErr: false,
		},
		{
			title: "bad rgb",
			cName: "rgb(999,33,216)",
			isErr: true,
		},
		{
			title: "pinkish",
			cName: "rgba(172,33,216,0.5)",
			isErr: false,
		},
		{
			title: "invalid alpha",
			cName: "rgba(172,33,216,55)",
			isErr: true,
		},
	} {

		l := LocalPenColour{}
		err := l.Unmarshal(test.cName)
		if err != nil && test.isErr == false {
			t.Errorf("%d test %s error %s", i, test.title, err)
		}
		t.Logf("test %d : title %s %+v", i, test.title, l)
	}
}
