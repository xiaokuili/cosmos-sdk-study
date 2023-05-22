package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

// add user:ignite account create id
// query bank: datafactoryd query bank balances cosmos1p00vhp362anpn3vt2wdp3d743v5y4ym99n3lvf
// query bank: datafactoryd query bank balances cosmos1sjyqpwfnhgv43flm6wdx883g4quvf2uwct7h29
// send tx: datafactoryd tx bank send cosmos1p00vhp362anpn3vt2wdp3d743v5y4ym99n3lvf cosmos189glceljqsmk5ahjnvhcvg82yjnemrljqqmnw4 100token
//

type Account struct {
	Key string `json:"key"`
}

type AccountResp struct {
	Key      string `json:"key"`
	Mnemonic string `json:"mnemonic"`
	Address  string `json:"address"`
}

type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Money    string `json:"money"`
}

type Bank struct {
	Address string `json:"address"`
}

type BankResp struct {
	Address string `json:"address"`
	Money   string `json:"money"`
}

var account Account
var transaction Transaction
var bank Bank

func Exec(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		return []byte(""), err
	}
	return stdout, nil
}

func check(r *http.Request) (string, bool) {
	if r.Method != http.MethodPost {
		return "Bad Request Method", true
	}
	return "", false
}

func ReqToStruct(r *http.Request, s interface{}) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("Read Error")
	}

	if err := json.Unmarshal(body, &s); err != nil {
		return errors.New("Json Error")
	}
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Welcome come!!! Happy In Blockchain DataFactory"))
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {

		if s, ok := check(r); ok {
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		err := ReqToStruct(r, &account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := account.Key
		if key == "" {
			http.Error(w, "请提供正确的key", http.StatusBadRequest)

			return
		}

		stdout, err := Exec("ignite", "account", "create", key)
		if err != nil {
			http.Error(w, "创建用户失败"+err.Error(), http.StatusBadRequest)
			return
		}

		re := regexp.MustCompile(`(?s)keep your mnemonic in a secret place:\n\n(.*)\n`)
		match := re.FindStringSubmatch(string(stdout))
		if len(match) < 2 {
			http.Error(w, "创建用户失败", http.StatusBadRequest)
			return
		}
		mnemonic := match[1]
		stdout, err = Exec("ignite", "account", "show", key)
		if err != nil {
			http.Error(w, "获取用户地址失败"+err.Error(), http.StatusBadRequest)
			return
		}

		re = regexp.MustCompile(`(?s).*?(cosmos.*/?) \tPubKeySecp256k1`)
		match = re.FindStringSubmatch(string(stdout))
		if len(match) < 2 {
			http.Error(w, "提取用户地址失败", http.StatusBadRequest)
			return
		}
		address := match[1]
		resp := AccountResp{
			Key:      key,
			Mnemonic: mnemonic,
			Address:  address,
		}
		var rst []byte
		if rst, err = json.Marshal(resp); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// 返回执行结果
		w.Write(rst)
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if s, ok := check(r); ok {
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		err := ReqToStruct(r, &transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = Exec("datafactoryd", "tx", "bank", "send", transaction.Sender, transaction.Receiver, transaction.Money, "-y")
		if err != nil {
			http.Error(w, "转账失败", http.StatusBadRequest)
			return
		}

		w.Write([]byte("succeed"))
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if s, ok := check(r); ok {
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		err := ReqToStruct(r, &bank)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rst, err := Exec("datafactoryd", "query", "bank", "balances", bank.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.Write(rst)
	})
	http.ListenAndServe(":8000", nil)
}
