#! /bin/bash
# By WJQSERVER-STUDIO_WJQSERVER
#https://github.com/WJQSERVER/tools-stable

# 导入配置文件
source "repo_url.conf"

mikublue="\033[38;2;57;197;187m"
yellow='\033[33m'
white='\033[0m'
green='\033[0;32m'
blue='\033[0;34m'
red='\033[31m'
gray='\e[37m'

#彩色
mikublue(){
    echo -e "\033[38;2;57;197;187m\033[01m$1\033[0m"
}
white(){
    echo -e "\033[0m\033[01m$1\033[0m"
}
red(){
    echo -e "\033[31m\033[01m$1\033[0m"
}
green(){
    echo -e "\033[32m\033[01m$1\033[0m"
}
yellow(){
    echo -e "\033[33m\033[01m$1\033[0m"
}
blue(){
    echo -e "\033[34m\033[01m$1\033[0m"
}
gray(){
    echo -e "\e[37m\033[01m$1\033[0m"
}
option(){
    echo -e "\033[32m\033[01m ${1}. \033[38;2;57;197;187m${2}\033[0m"
}

version=$(curl -s --max-time 3 ${repo_url}Version)
if [ $? -ne 0 ]; then
    version="unknown"  
fi

clear

function writeSystemdServiceNoEnv() {

cat <<EOF > /etc/systemd/system/caddydash.service
[Unit]
Description=CaddyDash
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=exec
User=root
Group=root
ExecStart=/root/data/caddy/caddydash -c ./config/config.toml
WorkingDirectory=/root/data/caddy
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target

EOF

}

function writeSystemdServiceWithPasswdEnv() {
    local username="$1"
    local password="$2"

cat <<EOF > /etc/systemd/system/caddydash.service
[Unit]
Description=CaddyDash
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=exec
User=root
Group=root
ExecStart=/root/data/caddy/caddydash -c ./config/config.toml
WorkingDirectory=/root/data/caddy
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE
Environment="CADDYDASH_USERNAME=${username}"
Environment="CADDYDASH_PASSWORD=${password}"

[Install]
WantedBy=multi-user.target

EOF

echo -e "${green}>${white} $mikublue CaddyDash 用户名: ${username}" $white
echo -e "${green}>${white} $mikublue CaddyDash 密码: ${password}" $white

}

# 显示免责声明
echo -e "${red}免责声明：${mikublue}请阅读并同意以下条款才能继续使用本脚本。"
echo -e "${yellow}===================================================================="
echo -e "${mikublue}本脚本仅供学习和参考使用，作者不对其完整性、准确性或实用性做出任何保证。"
echo -e "${mikublue}使用本脚本所造成的任何损失或损害，作者不承担任何责任。"
echo -e "${mikublue}不提供/保证任何功能的可用性，安全性，有效性，合法性"
echo -e "${mikublue}当前版本为${white}  [${yellow} V${version} ${white}]  ${white}"
echo -e "${yellow}===================================================================="
sleep 1

# 判断caddy是否已安装
if [ -f /root/data/caddy/caddy ]; then
    echo -e "[${yellow}Warning${white}] $mikublue caddy已安装, 将会覆盖安装" $white
    read -p "是否继续安装? [Y/n]" choice
    if [[ "$choice" == "" || "$choice" == "Y" || "$choice" == "y" ]]; then
        echo -e "[${yellow}Warning${white}] $mikublue 停止caddy" $white
        systemctl stop caddy
        systemctl stop caddydash
        systemctl disable caddy
        systemctl disable caddydash
        else
            echo -e "[${red}Cancel${white}] $mikublue 取消安装" $white
            exit 1
    fi
fi

echo -e "[${yellow}RUN${white}] $mikublue 開始安裝Caddy" $white

echo -e "${green}>${white} $mikublue 創建安裝目錄" $white
mkdir -p /root/data/caddy
mkdir -p /root/data/caddy/config.d
mkdir -p /root/data/caddy/config
echo -e "${green}>${white} $mikublue 下載主程序" $white
input_version="$@" #获取输入的版本号
if [ -z "$input_version" ]; then
    VERSION=$(curl -s https://raw.githubusercontent.com/WJQSERVER/caddy/main/TEST-VERSION)
else
    VERSION=$input_version
fi
wget -q -O /root/data/caddy/caddy.tar.gz https://github.com/WJQSERVER/caddy/releases/download/$VERSION/caddy-linux-amd64-pages.tar.gz
echo -e "${green}>${white} $mikublue 解壓程序及其資源" $white
tar -xzvf /root/data/caddy/caddy.tar.gz -C /root/data/caddy
echo -e "${green}>${white} $mikublue 清理安裝資源" $white
rm /root/data/caddy/caddy.tar.gz
echo -e "${green}>${white} $mikublue 設置程序運行權限" $white
chmod +x /root/data/caddy/caddy
chown root:root /root/data/caddy/caddy

# 询问是否创建caddydash账户密码
echo -e "${green}>${white} $mikublue 是否创建caddydash账户密码" $white

read -p "是否创建caddydash账户密码? [Y/n]" choice
if [[ "$choice" == "" || "$choice" == "Y" || "$choice" == "y" ]]; then
    read -p "请输入CaddyDash用户名: " caddydash_username
    read -s -p "请输入CaddyDash密码: " caddydash_password
    writeSystemdServiceWithPasswdEnv "$caddydash_username" "$caddydash_password"
else
    writeSystemdServiceNoEnv
fi

echo -e "${green}>${white} $mikublue 創建SERVICE文件" $white


echo -e "${green}>${white} $mikublue 拉取Caddyfile配置" $white
# 预留配置

echo -e "${green}>${white} $mikublue 啟動程序" $white
systemctl daemon-reload
systemctl enable caddydash.service
systemctl start caddydash.service
echo -e "[${green}OK${white}] $mikublue caddy安装完成" $white

#回到root目录
cd /root

# 导入配置文件
source "repo_url.conf"

#等待1s
sleep 1

#返回菜单/退出脚本
read -p "是否返回菜单?: [Y/n]" choice

if [[ "$choice" == "" || "$choice" == "Y" || "$choice" == "y" ]]; then
    wget -O program-menu.sh ${repo_url}program/program-menu.sh && chmod +x program-menu.sh && ./program-menu.sh
else
    echo "脚本结束"
fi
