package utils

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strconv"
	"wifidog-server/dao"
	//"golang.org/x/crypto/ssh/terminal"
	"log"
)

type SshConfig struct {
	Host    string
	Type    string
	KeyPath string
	Port    string
	config  ssh.ClientConfig
}

var apDao = new(dao.ApDao)

func ApConfig(sshhost, ssid, lan_ip string) error {
	sshconfig := SshInit(sshhost)

	//获取ssh client
	client := SshClient(sshconfig)
	defer client.Close()

	//创建ssh-session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()

	//远程命令
	cmd := TransToCommend(ssid, lan_ip)

	//执行命令
	err = ExecutAndCheck(client, cmd)
	if err != nil {
		return err
	} else {
		log.Println("AP配置成功")
		return nil
	}
}

//初始化ssh的配置
func SshInit(host string) SshConfig {
	var Conf SshConfig
	Conf.Host = host
	Conf.KeyPath = "id_rsa"
	Conf.Port = apDao.GetForwardPort(host)
	hostKey, _ := getHostKey(host, Conf.Port)

	//创建ssh登陆配置
	Conf.config = ssh.ClientConfig{
		//Timeout:         30*time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		Timeout:         0,
		User:            "root",
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		Auth:            []ssh.AuthMethod{publicKeyAuthFunc(Conf.KeyPath)},
	}
	return Conf
}

//创建ssh的client
func SshClient(conf SshConfig) *ssh.Client {
	port, _ := strconv.Atoi(conf.Port)
	addr := fmt.Sprintf("%s:%d", "localhost", port)
	sshClient, err := ssh.Dial("tcp", addr, &(conf.config))
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	} else {
		log.Println("创建ssh client 成功")
	}
	return sshClient
}

//配置参数翻译为配置命令行
func TransToCommend(ssid, lan_ip string) string {
	//cmd := "uci set wireless.default_radio0.ssid='" + ssid + "';uci commit;wifi;sleep 5;" +
	//	"uci set network.lan.ipaddr='" + lan_ip + "';uci commit;/etc/init.d/network restart;sleep 5"
	cmd := "uci set wireless.default_radio0.ssid=" + ssid + ";uci commit;wifi;sleep 5;uci set network.lan.ipaddr='" + lan_ip + "';uci commit;/etc/init.d/network restart;sleep 5"
	log.Printf("execute the cmd:%s", cmd)
	return cmd
}

//配置命令的执行与结果验证
func ExecutAndCheck(client *ssh.Client, cmd string) error {
	//创建ssh-session
	session, err := client.NewSession()
	if err != nil {
		//log.Fatal("创建ssh session 失败", err)
		log.Println(err)
	} else {
		log.Println("create ssh session success")
	}
	//方式1：
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		//log.Fatal("远程执行cmd 失败", err)
		log.Printf("cmd: %s, failed", string(combo))
		return err
	} else {
		log.Printf("cmd：%s ,success", string(combo))
	}
	return nil
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	//keyPath, err := homedir.Expand(kPath)
	log.Printf("keypath:%s", kPath)
	key, err := ioutil.ReadFile(kPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

//从文本地中获取中获取要远程登录主机的hostkey
func getHostKey(host, port string) (ssh.PublicKey, error) {
	var hostKey ssh.PublicKey
	db_hostkey := apDao.GetApHostkeyDB(host, port)
	//在本地数据库中获取ssh连接的服务器端的公钥文件夹，并将其转换为ssh.PublicKey格式
	var err error
	// ParseAuthorizedKeys parses a public key from an authorized_keys file used in OpenSSH
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(db_hostkey))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing: %v", err))
	}
	hostKey, err = ssh.ParsePublicKey(pk.Marshal())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing: %v", err))
	}

	if hostKey == nil {
		return nil, errors.New(fmt.Sprintf("no hostkey for %s", host))
	}
	return hostKey, nil
}
