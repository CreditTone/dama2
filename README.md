打码兔网页版功能说明
简介
打码兔网页版提供预授权、登陆、注册新用户、查询信息（查询用户信息、查询余额）、打码（通过URL打码、POST图片打码、查询打码结果、报告错误）等功能。
除POST图片打码仅支持POST方法外，其他功能都支持POST和GET两种方法。
查询信息和打码前需执行“预授权”请求，以“预授权”返回的预授权信息执行“登陆”请求，成功登陆后返回授权信息，查询和打码请求需附带授权信息。所有请求在成功时都会返回授权信息，该授权信息可用于后续的其他请求。
注册新用户前先执行“预授权”请求，以“预授权”返回的预授权信息执行“注册新用户请求”。
请求参数
请求参数包括：



响应数据格式
所有请求的结果数据均以JSON格式返回到客户端。
响应数据编码为“UTF-8”。
请求失败时仅返回“ret”和“desc”两个字段。如下所示为打码请求失败的响应数据（加粗部分因不同的错误而有所不同）：
{"ret":"-304","desc":"余额不足，题分不足"}

ret指返回码，失败时均为负数； desc为错误描述。

成功时所有请求都会返回“ret”、“desc”、“auth”；打码请求还会返回“id”，用于查询结果；查询打码结果请求还会返回“result”；查询余额请求会返回“balance”；查询用户信息会返回“qq”、“email”、“tel”。
查询余额成功的例子响应数据如下：
{"ret":"0","desc":"成功","balance":"0","auth":"1382614466607AuuSuzazOsRziqvnD9NWCKJQzbMiB52bXJhstwX8O1vnoh%2FP5SBzLUK0o9o1%2BXXyDfe6lyWdnPsoT%0D%0AFefql2B9aw%3D%3D"}
ret指返回码，成功时为0； balance为余额； auth为返回的授权信息，作为后续调用的参数。
响应字段


功能介绍
预授权
功能URL：http://api.dama2.com:7788/app/preauth
输入参数：无
返回数据：ret，desc，auth(预授信息，仅用于计算加密信息)
登陆
功能URL：http://api.dama2.com:7788/app/login
输入参数：appID、sname(可选参数)、encinfo
返回数据：ret，desc，auth
注册新用户
功能URL：http://api.dama2.com:7788/app/register
输入参数：appID、sname(可选参数)、encinfo、qq、email、tel
返回数据：ret，desc，auth
查询用户信息
功能URL：http://api.dama2.com:7788/app/readInfo
输入参数：auth
返回数据：ret，desc，auth, name, qq、email、tel
查询余额
功能URL：http://api.dama2.com:7788/app/getBalance
输入参数：auth
返回数据：ret，desc，auth
POST文件打码
功能URL：http://api.dama2.com:7788/app/decode
输入参数：auth、type、len(可选)、timeout（可选）、文件数据
返回数据：ret，desc，aut，id(用于查询结果和报告结果)
通过URL打码
功能URL：http://api.dama2.com:7788/app/decodeURL
输入参数：auth、type、len(可选)、timeout（可选）、url、cookie（可选）、referer（可选）
返回数据：ret，desc，aut，id（用于查询结果和报告结果）

查询打码结果  
功能URL：http://api.dama2.com:7788/app/getResult
输入参数：auth，id
返回数据：ret，desc，result，cookie, auth

报告结果  
功能URL：http://api.dama2.com:7788/app/reportError
输入参数：auth，id
返回数据：ret，desc，auth
附录
返回码
客户端程序必须特别关注的特殊返回码如下表所示：


完整的返回码如下表所示：

加密信息的计算算法
伪代码如下（假设软件KEY为“9503ce045ad14d83ea876ab578bd3184”， 预授权字串为“1234567890abcdef1234567890abcdef”，用户名为“user”，密码为“test”，结果为“a733506fda6e182300d34a1bcca568d8a733506fda6e182300d34a1bcca568d8f6368183d9f688b7a43c783d0184850751f1d013925667735ddd457beaa8da8010b74e6f64b15c4d”）：





Java代码如下：
private String byteArray2HexString(byte [] data) {
    StringBuilder sb = new StringBuilder();
for (byte b : data) {
        String s = Integer.toHexString(b & 0xff);
        if (s.length() == 1) {
            sb.append("0" + s);
        } else {
            sb.append(s);
        }
    }
    return sb.toString();
}

String password = "test";   //用户密码明文
//压缩软件KEY为8字节，用作DES加密的KEY
byte [] key16 = hexString2ByteArray("9503ce045ad14d83ea876ab578bd3184"); 
byte [] key8 = new byte[8];
for (int i = 0; i < 8; i++) {
    key8[i] = (byte)((key16[i] ^ key16[i + 8]) & 0xff);
}

String result = null;
byte [] data = password.getBytes();
try {
    java.security.MessageDigest md5 = java.security.MessageDigest.getInstance("MD5");
    md5.update(data, start, len);
    data = md5.digest();
pwd_md5_str = byteArray2HexString(data); //转为16进制字符串

String enc_data_str = preauth + “\n” + user_name + “\n” + pwd_md5_str;

SecureRandom sr = new SecureRandom();
DESKeySpec dks = new DESKeySpec(key);
SecretKeyFactory keyFactory = SecretKeyFactory.getInstance("DES");
SecretKey securekey = keyFactory.generateSecret(dks);
Cipher cipher = Cipher.getInstance("DES"); 
cipher.init(Cipher.ENCRYPT_MODE , key8, sr);
byte [] resultData = cipher.doFinal(enc_data_str.getBytes());
result = byteArray2HexString(resultData);
} catch (NoSuchAlgorithmException e) {
    e.printStackTrace();
} 