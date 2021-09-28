/***
export const formatTagName = (name: string): string => {
  const invalidRegex = /[\/\_]/gi;
  const noInvalidChars = name.replace(invalidRegex, '-');
  return noInvalidChars.replace(/-+$/, ''); // replace trailing '-' characters too
};
*/

package format

import (
	"testing"
)

func TestFormatBranch(t *testing.T) {
	type args struct {
		branch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid branch",
			args: args{
				branch: "good-branch",
			},
			want: "good-branch",
		},
		{
			name: "valid branch with slash swapped",
			args: args{
				branch: "good/branch",
			},
			want: "good-branch",
		},
		{
			name: "valid branch with dash ended",
			args: args{
				branch: "good-branch-",
			},
			want: "good-branch",
		},
		{
			name: "valid branch with dashes ended",
			args: args{
				branch: "good-branch---",
			},
			want: "good-branch",
		},
		{
			name: "valid branch",
			args: args{
				branch: "good////--branch",
			},
			want: "good------branch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBranch(tt.args.branch); got != tt.want {
				t.Errorf("FormatBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortBranchName(t *testing.T) {
	type args struct {
		branch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid branch",
			args: args{
				branch: "good-branch",
			},
			want: "good-branch",
		},
		{
			name: "valid branch long",
			args: args{
				branch: "good-branch-this-is-a-very-long-branch-that-is-around-a-lot-of-characters",
			},
			want: "good-branch-this-is-a-very-long-branch-that-is-around-a-lot-of-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortBranchName(tt.args.branch); got != tt.want {
				t.Errorf("ShortBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}
