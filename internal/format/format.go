/***
export const formatTagName = (name: string): string => {
  const invalidRegex = /[\/\_]/gi;
  const noInvalidChars = name.replace(invalidRegex, '-');
  return noInvalidChars.replace(/-+$/, ''); // replace trailing '-' characters too
};
*/

package format

import (
	"fmt"
	"regexp"
)

const (
	MaxLabelLength = 63
)

// FormatBranch will modify the branch name input to something that can be accepted as a label value.
func FormatBranch(branch string) string {
	re1 := regexp.MustCompile(`[\/\_]`) // remove typical underscores / slashes
	a := re1.ReplaceAll([]byte(branch), []byte("-"))
	re2 := regexp.MustCompile(`-+$`) // replace trailing '-' characters too
	return fmt.Sprintf("%s", re2.ReplaceAll(a, []byte("")))
}

// ShortBranchName is to limit the length of the branch string that is searched for
func ShortBranchName(branch string) string {
	if len(branch) > MaxLabelLength {
		return branch[:MaxLabelLength]
	}
	return branch
}
