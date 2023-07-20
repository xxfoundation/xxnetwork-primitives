////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package notifications

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// Tests that a list of Data CSV encoded by BuildNotificationCSV and decoded bu
// DecodeNotificationsCSV matches the original.
func TestBuildNotificationCSV_DecodeNotificationsCSV(t *testing.T) {
	rng := rand.New(rand.NewSource(186745))
	expected := make([]*Data, 50)
	for i := range expected {
		identityFP, messageHash := make([]byte, 25), make([]byte, 32)
		rng.Read(messageHash)
		rng.Read(identityFP)
		expected[i] = &Data{IdentityFP: identityFP, MessageHash: messageHash}
	}

	csvData, _ := BuildNotificationCSV(expected, 9999)
	dataList, err := DecodeNotificationsCSV(string(csvData))
	if err != nil {
		t.Errorf("Failed to decode notifications CSV: %+v", err)
	}

	if !reflect.DeepEqual(expected, dataList) {
		t.Errorf("The generated Data list does not match the original."+
			"\nexpected: %v\nreceived: %v", expected, dataList)
	}
}

// Consistency test of BuildNotificationCSV.
func TestBuildNotificationCSV(t *testing.T) {
	expected := `U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==
GsvgcJsHWAg/YdN1vAK0HfT5GSnhj9qeb4LlTnSOgec=,nku9b+NM3LqEPujWPoxP/hzr6lRtj6wT3Q==
GqwEzi6ih3xVec+ix44bC6+uiBuCp1EQikLtPJA8qkM=,Rlp4YgYWl4rtDOPGxPOue8PgSVtXEv79vg==
DBAoh+EA2s0tiF9pLLYH2gChHBxwceeWotwtwlpbdLI=,4SlwXic/BckjJoKOKwVuOBdljhBhSYlH/Q==
80RBDtTBFgI/qONXa2/tJ/+JdLrAyv2a0FaSsTYZ5zg=,lk39x56NU0NzZhz9ZtdP7B4biUkatyNuSw==
dSFikM8r60LDyicyhWDxqsBnzqbov0bUqytGgEAsX7I=,gg6IXTJg8d6XgoPUoJo2+WwglBdG4+1Npg==
Rqmui0+ntPw6ILr6GnXtMnqGuLDDmvHP0rO1EhnqeVM=,Or9EjSxHnTJgdTOQWRTIzBzwnaOeDpKdAg==
Sry8sWk5e7c05+8KbgHxhU3rX+Qk/vesIQiR9ZdeKSo=,oriqBHxhzbMzc+vnLCegmMAhl9rmtzLDUQ==
32aPh04snxzgnKhgF+fiF0gwP/QcGyPhHEjtF1OdaF8=,dvKnmLxk3g5dsoZLKtPCbOY4I0J2WhPWlg==
5S33YPbDRl4poNykasOg1XATO8IVcfX1SmQxBVE/2EI=,mxlK4bqfKoOGrnKzZh/oLCrGTb9GFRgk4g==
MFMSY3yZwrh9bfDdXvKDZxkHLWcvYfqgvob0V5Iew3w=,DkYM8NcD0H3F9WYaRQEzQJpxK2pmq9e6ZQ==
IkyiaXjZpc5i/rEag48WYi61TO4+Z1UinBg8GTOpFlg=,Xhg7twkZLbDmyNcJudc4O5k8aUmZRbCwzw==
49wuwfyWENfusZ0JFqJ0I8KeRC8OMcLJU5Zg8F+zfkU=,zRvwvPwaNGxDTxHPAEFvphaVuSAuaDY6HA==
eH9HhOCu2ceFZBhOEx8efIEfvYhbzGc06JM/PLLyXVI=,+fjHVHrX4dYnjJ98hy+ED52U2f3trpPbJA==
lXGPWAuMjyvsxqp2w7D5SK++YSelz9VrwRs8Lqg3ocY=,aagi92hk7CrgzWv93yGxFER0v9N80ga1Gg==
zgUKthmex7OW1hj94OGimZpvPZ+LergUn3Leulxs1P0=,TTkskrSyGsgSA0Bi38MGOnpoYrD+8QUpGQ==
wqh6SZT8HkAeEcWknH6OqeZdbMQEZf01LyxC7D0+9g0=,tpdAUX3HZSue7/UWU1qhyfM9sT7R964b4w==
hBMjKNED+HGvm80VIzw5OXj1wXCJ6PMmegzMfjm/ysc=,rEK+LBcsYkPRBjMDbT1GuBkWrkb/E9amsg==
+tkHnW3zRAWQKWZ7LrQaQAEXVW/ly0BbMXCsrKXHW68=,f3tw6GFE07oDazsfWP5CeVDn0E9MJuvhLw==
1eEjcZgIogS4Ubps8spsRu2IFi9dRc21oHY65+GDP7c=,rfmJNvUeTdqKKE7xmoW7h0N7QQMwWjs4bA==
fTbZLSJUmWCnFPKoKeHCAhZzvzDFC2edUFaJVcnBmAg=,nZX2A5fSr1+PyREL46nhJelEhJeXCNaqfA==
/GKejZFHzy9ftqBVkauGhzoerQWkpmcdaVFcg53Yrzo=,Otd0AsX9OoOgRgipiTMAIWLdTB/1VH9XUg==
Ax8hIeFBCKpaV0VsrpHBcymtWs5h6um2Ut8zALTCq1g=,J3bYW2jKMtXDc8JkeFg7xI+ja+SNZZw/4Q==
c0EBx+SP03+5+uPwu06bbfR1Ki6RZM8F9WjSyJ6k1l0=,dgYOZIeQWTJLt1rbFBovfC/eeBH0gc8Iag==
PsPYs3cAEv0npLZbAq6FJW9zbt4+TdhXIJV1pIjVdA0=,L3JpWlcNvyZH8pXiM5Xu2s/2NuGwzyDeag==
EP+ZQ+3Kb5a/TdrwC51PzWrL27P2MZRQNYaopliuYLU=,7lOata0Z8roj3KZn36ZVE0xZSiyAa9+k5w==
VuRYtIuuSQ0ELgejVels+4nMq/KBnXlNnhKC/QpyVPE=,s5T26rxmpki639tH01CKaTgLpg1f9LQyew==
9lDgExuPV3WthpenNGPNKAbmru75K16b/+QOlGaZD6M=,rsEeSny2rrsXt/7SlRPTHtT/HRbm1ZlWGQ==
UV9fU4dpAO0PetHyOLszRAnjwWSVc6VvQ6jh0hNyRvw=,psYzJNQ/g+wNTS/WUG/f7uIeJDI9gOfLhA==
X0PSIyKapCEUSifbt8RAwceY+aJNLIXxLCSIv4fS2Sk=,oKR5pVt7c+TvskFDTjbUT315OI2hnlz+gw==
IWN1mCbfOfgzaVyiKqZRlUiQvNzPZq09c6jhq5+Dh30=,Jju7J/W8SXvWVEdNy4YqtN1om6BNDa5ooA==
g9al6HTEHOSudp3dtiHBZDI5vTeKLpGprOJ38sCNcUs=,ydnmLAViyiluEqd2F0TduCOoLxm6fQpSSw==
VJK79yrDTvy5Cl7fbbwhn78w7PJfpmbJJGsIHV0sV44=,7wAIsI1hoJdkBPQuqCpIc/sNZId3faZHBg==
t0hXpZ8dKn82F6O81VqVn9GSBMLjv6zg5gMLfABpuXw=,aQyYNMIoKbqb1P+pr1gZb3deMPPJO0nsLw==
Hoo35EiId/9sfiQlH+FQ7OfMWvss7XprvKzj7qAb11k=,QA2HuYCzU8EVy8Fip3jdnqBCNZ1MIP4hig==
Rm/cqgfzRclH5aCWoj+JZ89P4Si96pz8xljy7bEkkpw=,3M9Yj0lOvjNGwZrteHuXxXcN/t6EXPWwQA==
3LYIlEhmP8MyF8HyL7TKpWBOFiDDl7Oo40e3k0PkPl4=,lPyl5AhHBG352IgCviQSoTRntmVWLzKHSA==
5IPF6phRI8xCLk96jOl0B1OPYfZ+ga42GtW89w8iiDE=,aw4ukENMK3yiyg2KICMlx7gMtjXoXb0jNw==
QNWTeKISlTt5F8x/RdbsAU0fC1kNaLRRMzwAisvlEjE=,+4CfIcugABlRkeMY0PNJ84IlHeA7NfV9zw==
UrloJgqUXJGcj7n7jfqEfWb7oCNK27w240akwcvimRg=,FGu6CxanGNbECj5ZdsoEaGR0nEgEx5zJrQ==
ZLZ2Bw9hP9+WSKJW3DwiOkvOiRWUK9lrAHMdrZWDfD8=,r/8aTMECHltCu3V4nHCj932lPCXgSLkLqg==
HrARGizMUEcrKECJa840U6mtBJct5H/GZEahdvtaE8I=,Xcu6Vrv2NV4bKvhmVDH3RyqWYYFmnxAfWg==
Vyy0GiAUFyBexvVbintbSsYQjuBFVTHkOGRH9fTJGdw=,S77jKfBIvvwO5ArLSmxuEHLQQwBQjdXzWw==
LPwGgdnQAZaEWYyCdG1Zk/AB99k9z/INedKtTv1e5Ow=,qyjyubYZBFj+NsS3dayvYMFUI5W2jO9WjQ==
OWA4Tr2KTqoq6+xmTlY4cNuAPSgOPmJwo7D+A4vILZw=,gw/oRNJWsLXpYvMxM58T2FKXOynKoD6QFA==
qIfiAe4BFutxC8au4sJOXZBExUpNymRkA2w2FMafnII=,PFvyIccm6amL8jQBONIh2lPeVMi1Bvk/fg==
AcsU15TF3uaMZzKcHTyptNP7EBq5eBYhI2vBK/rFKCQ=,Gcam+D1Hzebx9Zs8AHd3yAALcOHAyJAiuQ==
2xNm0x0FAN2fAkPW6rUP0gFhx0hJw94sUaubeM+WWRA=,iC3H9TvHMgsc9IRy9ks2Qd/TaY9zTNkOXA==
A3hMWMAcrvqWoVNZPxQqYFWLMoCUCnrl2NArseYXnTk=,WsPBzNwVH8QF0fcpHDoq7po6JHhgL9Zcew==
`
	extra := "Zq3/Nor7+NgAzkvg7LxVOYyRMMnAEDxkHpGnKpeHltc=,wGc+G+CLk/qEIoGMQ0XBZlyHkiYS3r7nkw==\n"

	rng := rand.New(rand.NewSource(42))
	dataList := make([]*Data, 50)
	for i := range dataList {
		identityFP, messageHash := make([]byte, 25), make([]byte, 32)
		rng.Read(messageHash)
		rng.Read(identityFP)
		dataList[i] = &Data{IdentityFP: identityFP, MessageHash: messageHash}
	}

	csv, rest := BuildNotificationCSV(dataList, 4096)
	second, _ := BuildNotificationCSV(rest, 4096)
	if expected != string(csv) {
		t.Errorf("First pass mismatch.\nexpected:\n[%s]\n\nreceived:\n[%s]",
			expected, string(csv))
	}
	if extra != string(second) {
		t.Errorf("Second pass mismatch.\nexpected:\n[%s]\n\nreceived:\n[%s]",
			extra, string(second))
	}
}

func TestBuildNotificationCSV_small(t *testing.T) {
	expected := `U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==
GsvgcJsHWAg/YdN1vAK0HfT5GSnhj9qeb4LlTnSOgec=,nku9b+NM3LqEPujWPoxP/hzr6lRtj6wT3Q==
`
	rng := rand.New(rand.NewSource(42))
	dataList := make([]*Data, 2)
	for i := range dataList {
		identityFP, messageHash := make([]byte, 25), make([]byte, 32)
		rng.Read(messageHash)
		rng.Read(identityFP)
		dataList[i] = &Data{IdentityFP: identityFP, MessageHash: messageHash}
	}

	csv, rest := BuildNotificationCSV(dataList, 4096)
	if expected != string(csv) {
		t.Errorf("First pass mismatch.\nexpected:\n[%s]\n\nreceived:\n[%s]",
			expected, string(csv))
	}
	if len(rest) != 0 {
		t.Errorf("Should not have been any overflow, but got %+v", rest)
	}
}

// Error path: Tests that DecodeNotificationsCSV returns the expected error for
// an invalid MessageHash.
func TestDecodeNotificationsCSV_InvalidMessageHashError(t *testing.T) {
	invalidCSV := `U4x/lrFkvxuXu59LtHLonnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==
`
	expectedErr := "Failed decode an element"
	_, err := DecodeNotificationsCSV(invalidCSV)
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Unexpected error for invalid MessageHash."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: Tests that DecodeNotificationsCSV returns the expected error for
// an invalid identityFP.
func TestDecodeNotificationsCSV_InvalididentityFPError(t *testing.T) {
	invalidCSV := `U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ1hRMiIU1hA==
`
	expectedErr := "Failed decode an element"
	_, err := DecodeNotificationsCSV(invalidCSV)
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Unexpected error for invalid identityFP."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: Tests that DecodeNotificationsCSV returns the expected error for
// an invalid identityFP.
func TestDecodeNotificationsCSV_NoEofError(t *testing.T) {
	invalidCSV := `U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==,"`
	expectedErr := "Failed to decode notifications CSV"
	_, err := DecodeNotificationsCSV(invalidCSV)
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Unexpected error for invalid identityFP."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}
