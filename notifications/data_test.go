package notifications

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestMake_DecodeNotificationsCSV(t *testing.T) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const numNotifications = 45

	notifList := make([]*Data, 0, numNotifications)

	for i := 0; i < numNotifications; i++ {
		msgHash := make([]byte, 32)
		ifp := make([]byte, 25)
		rng.Read(msgHash)
		rng.Read(ifp)
		notifList = append(notifList, &Data{MessageHash: msgHash, IdentityFP: ifp})
	}

	notifCSV, _ := BuildNotificationCSV(notifList, 4096)
	newNotifList, err := DecodeNotificationsCSV(string(notifCSV))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(notifList, newNotifList) {
		t.Errorf("The generated notifivations do not match")
	}
}

func TestBuildNotificationCSV(t *testing.T) {
	expected := "U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==" +
		"\nGsvgcJsHWAg/YdN1vAK0HfT5GSnhj9qeb4LlTnSOgec=,nku9b+NM3LqEPujWPoxP/hzr6lRtj6wT3Q==" +
		"\nGqwEzi6ih3xVec+ix44bC6+uiBuCp1EQikLtPJA8qkM=,Rlp4YgYWl4rtDOPGxPOue8PgSVtXEv79vg==" +
		"\nDBAoh+EA2s0tiF9pLLYH2gChHBxwceeWotwtwlpbdLI=,4SlwXic/BckjJoKOKwVuOBdljhBhSYlH/Q==" +
		"\n80RBDtTBFgI/qONXa2/tJ/+JdLrAyv2a0FaSsTYZ5zg=,lk39x56NU0NzZhz9ZtdP7B4biUkatyNuSw==" +
		"\ndSFikM8r60LDyicyhWDxqsBnzqbov0bUqytGgEAsX7I=,gg6IXTJg8d6XgoPUoJo2+WwglBdG4+1Npg==" +
		"\nRqmui0+ntPw6ILr6GnXtMnqGuLDDmvHP0rO1EhnqeVM=,Or9EjSxHnTJgdTOQWRTIzBzwnaOeDpKdAg==" +
		"\nSry8sWk5e7c05+8KbgHxhU3rX+Qk/vesIQiR9ZdeKSo=,oriqBHxhzbMzc+vnLCegmMAhl9rmtzLDUQ==" +
		"\n32aPh04snxzgnKhgF+fiF0gwP/QcGyPhHEjtF1OdaF8=,dvKnmLxk3g5dsoZLKtPCbOY4I0J2WhPWlg==" +
		"\n5S33YPbDRl4poNykasOg1XATO8IVcfX1SmQxBVE/2EI=,mxlK4bqfKoOGrnKzZh/oLCrGTb9GFRgk4g==" +
		"\nMFMSY3yZwrh9bfDdXvKDZxkHLWcvYfqgvob0V5Iew3w=,DkYM8NcD0H3F9WYaRQEzQJpxK2pmq9e6ZQ==" +
		"\nIkyiaXjZpc5i/rEag48WYi61TO4+Z1UinBg8GTOpFlg=,Xhg7twkZLbDmyNcJudc4O5k8aUmZRbCwzw==" +
		"\n49wuwfyWENfusZ0JFqJ0I8KeRC8OMcLJU5Zg8F+zfkU=,zRvwvPwaNGxDTxHPAEFvphaVuSAuaDY6HA==" +
		"\neH9HhOCu2ceFZBhOEx8efIEfvYhbzGc06JM/PLLyXVI=,+fjHVHrX4dYnjJ98hy+ED52U2f3trpPbJA==" +
		"\nlXGPWAuMjyvsxqp2w7D5SK++YSelz9VrwRs8Lqg3ocY=,aagi92hk7CrgzWv93yGxFER0v9N80ga1Gg==" +
		"\nzgUKthmex7OW1hj94OGimZpvPZ+LergUn3Leulxs1P0=,TTkskrSyGsgSA0Bi38MGOnpoYrD+8QUpGQ==" +
		"\nwqh6SZT8HkAeEcWknH6OqeZdbMQEZf01LyxC7D0+9g0=,tpdAUX3HZSue7/UWU1qhyfM9sT7R964b4w==" +
		"\nhBMjKNED+HGvm80VIzw5OXj1wXCJ6PMmegzMfjm/ysc=,rEK+LBcsYkPRBjMDbT1GuBkWrkb/E9amsg==" +
		"\n+tkHnW3zRAWQKWZ7LrQaQAEXVW/ly0BbMXCsrKXHW68=,f3tw6GFE07oDazsfWP5CeVDn0E9MJuvhLw==" +
		"\n1eEjcZgIogS4Ubps8spsRu2IFi9dRc21oHY65+GDP7c=,rfmJNvUeTdqKKE7xmoW7h0N7QQMwWjs4bA==" +
		"\nfTbZLSJUmWCnFPKoKeHCAhZzvzDFC2edUFaJVcnBmAg=,nZX2A5fSr1+PyREL46nhJelEhJeXCNaqfA==" +
		"\n/GKejZFHzy9ftqBVkauGhzoerQWkpmcdaVFcg53Yrzo=,Otd0AsX9OoOgRgipiTMAIWLdTB/1VH9XUg==" +
		"\nAx8hIeFBCKpaV0VsrpHBcymtWs5h6um2Ut8zALTCq1g=,J3bYW2jKMtXDc8JkeFg7xI+ja+SNZZw/4Q==" +
		"\nc0EBx+SP03+5+uPwu06bbfR1Ki6RZM8F9WjSyJ6k1l0=,dgYOZIeQWTJLt1rbFBovfC/eeBH0gc8Iag==" +
		"\nPsPYs3cAEv0npLZbAq6FJW9zbt4+TdhXIJV1pIjVdA0=,L3JpWlcNvyZH8pXiM5Xu2s/2NuGwzyDeag==" +
		"\nEP+ZQ+3Kb5a/TdrwC51PzWrL27P2MZRQNYaopliuYLU=,7lOata0Z8roj3KZn36ZVE0xZSiyAa9+k5w==" +
		"\nVuRYtIuuSQ0ELgejVels+4nMq/KBnXlNnhKC/QpyVPE=,s5T26rxmpki639tH01CKaTgLpg1f9LQyew==" +
		"\n9lDgExuPV3WthpenNGPNKAbmru75K16b/+QOlGaZD6M=,rsEeSny2rrsXt/7SlRPTHtT/HRbm1ZlWGQ==" +
		"\nUV9fU4dpAO0PetHyOLszRAnjwWSVc6VvQ6jh0hNyRvw=,psYzJNQ/g+wNTS/WUG/f7uIeJDI9gOfLhA==" +
		"\nX0PSIyKapCEUSifbt8RAwceY+aJNLIXxLCSIv4fS2Sk=,oKR5pVt7c+TvskFDTjbUT315OI2hnlz+gw==" +
		"\nIWN1mCbfOfgzaVyiKqZRlUiQvNzPZq09c6jhq5+Dh30=,Jju7J/W8SXvWVEdNy4YqtN1om6BNDa5ooA==" +
		"\ng9al6HTEHOSudp3dtiHBZDI5vTeKLpGprOJ38sCNcUs=,ydnmLAViyiluEqd2F0TduCOoLxm6fQpSSw==" +
		"\nVJK79yrDTvy5Cl7fbbwhn78w7PJfpmbJJGsIHV0sV44=,7wAIsI1hoJdkBPQuqCpIc/sNZId3faZHBg==" +
		"\nt0hXpZ8dKn82F6O81VqVn9GSBMLjv6zg5gMLfABpuXw=,aQyYNMIoKbqb1P+pr1gZb3deMPPJO0nsLw==" +
		"\nHoo35EiId/9sfiQlH+FQ7OfMWvss7XprvKzj7qAb11k=,QA2HuYCzU8EVy8Fip3jdnqBCNZ1MIP4hig==" +
		"\nRm/cqgfzRclH5aCWoj+JZ89P4Si96pz8xljy7bEkkpw=,3M9Yj0lOvjNGwZrteHuXxXcN/t6EXPWwQA==" +
		"\n3LYIlEhmP8MyF8HyL7TKpWBOFiDDl7Oo40e3k0PkPl4=,lPyl5AhHBG352IgCviQSoTRntmVWLzKHSA==" +
		"\n5IPF6phRI8xCLk96jOl0B1OPYfZ+ga42GtW89w8iiDE=,aw4ukENMK3yiyg2KICMlx7gMtjXoXb0jNw==" +
		"\nQNWTeKISlTt5F8x/RdbsAU0fC1kNaLRRMzwAisvlEjE=,+4CfIcugABlRkeMY0PNJ84IlHeA7NfV9zw==" +
		"\nUrloJgqUXJGcj7n7jfqEfWb7oCNK27w240akwcvimRg=,FGu6CxanGNbECj5ZdsoEaGR0nEgEx5zJrQ==" +
		"\nZLZ2Bw9hP9+WSKJW3DwiOkvOiRWUK9lrAHMdrZWDfD8=,r/8aTMECHltCu3V4nHCj932lPCXgSLkLqg==" +
		"\nHrARGizMUEcrKECJa840U6mtBJct5H/GZEahdvtaE8I=,Xcu6Vrv2NV4bKvhmVDH3RyqWYYFmnxAfWg==" +
		"\nVyy0GiAUFyBexvVbintbSsYQjuBFVTHkOGRH9fTJGdw=,S77jKfBIvvwO5ArLSmxuEHLQQwBQjdXzWw==" +
		"\nLPwGgdnQAZaEWYyCdG1Zk/AB99k9z/INedKtTv1e5Ow=,qyjyubYZBFj+NsS3dayvYMFUI5W2jO9WjQ==" +
		"\nOWA4Tr2KTqoq6+xmTlY4cNuAPSgOPmJwo7D+A4vILZw=,gw/oRNJWsLXpYvMxM58T2FKXOynKoD6QFA==" +
		"\nqIfiAe4BFutxC8au4sJOXZBExUpNymRkA2w2FMafnII=,PFvyIccm6amL8jQBONIh2lPeVMi1Bvk/fg==" +
		"\nAcsU15TF3uaMZzKcHTyptNP7EBq5eBYhI2vBK/rFKCQ=,Gcam+D1Hzebx9Zs8AHd3yAALcOHAyJAiuQ==" +
		"\n2xNm0x0FAN2fAkPW6rUP0gFhx0hJw94sUaubeM+WWRA=,iC3H9TvHMgsc9IRy9ks2Qd/TaY9zTNkOXA==" +
		"\nA3hMWMAcrvqWoVNZPxQqYFWLMoCUCnrl2NArseYXnTk=,WsPBzNwVH8QF0fcpHDoq7po6JHhgL9Zcew==\n"
	extra := "Zq3/Nor7+NgAzkvg7LxVOYyRMMnAEDxkHpGnKpeHltc=,wGc+G+CLk/qEIoGMQ0XBZlyHkiYS3r7nkw==\n"

	rng := rand.New(rand.NewSource(42))

	const numNotifications = 50

	notifList := make([]*Data, 0, numNotifications)

	for i := 0; i < numNotifications; i++ {
		msgHash := make([]byte, 32)
		ifp := make([]byte, 25)
		rng.Read(msgHash)
		rng.Read(ifp)
		notifList = append(notifList, &Data{MessageHash: msgHash, IdentityFP: ifp})
	}

	csv, rest := BuildNotificationCSV(notifList, 4096)

	second, _ := BuildNotificationCSV(rest, 4096)

	if expected != string(csv) {
		t.Errorf("First pass mismatch: Expected:\n[%s]\n, received: \n[%s]\n", expected, string(csv))
	}
	if extra != string(second) {
		t.Errorf("Second pass mismatch: Expected: \n[%s]\n, received: \n[%s]\n", extra, string(second))
	}
}

func TestBuildNotificationCSV_small(t *testing.T) {
	expected := "U4x/lrFkvxuXu59LtHLon1sUhPJSCcnZND6SugndnVI=,39ebTXZCm2F6DJ+fDTulWwzA1hRMiIU1hA==" +
		"\nGsvgcJsHWAg/YdN1vAK0HfT5GSnhj9qeb4LlTnSOgec=,nku9b+NM3LqEPujWPoxP/hzr6lRtj6wT3Q==\n"

	rng := rand.New(rand.NewSource(42))

	const numNotifications = 2

	notifList := make([]*Data, 0, numNotifications)

	for i := 0; i < numNotifications; i++ {
		msgHash := make([]byte, 32)
		ifp := make([]byte, 25)
		rng.Read(msgHash)
		rng.Read(ifp)
		notifList = append(notifList, &Data{MessageHash: msgHash, IdentityFP: ifp})
	}

	csv, rest := BuildNotificationCSV(notifList, 4096)

	if expected != string(csv) {
		t.Errorf("First pass mismatch: Expected:\n[%s]\n, received: \n[%s]\n", expected, string(csv))
	}
	if len(rest) != 0 {
		t.Errorf("Should not have been any overflow, but got [%+v]", rest)
	}
}
