#!/bin/bash

# 数据字典初始化脚本
# 用途：初始化系统所需的各类数据字典

# 加载配置
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
CONFIG_FILE="$SCRIPT_DIR/../conf/config.yaml"

# 从 config.yaml 读取数据库配置
DB_HOST=$(grep -A 15 "mysql:" "$CONFIG_FILE" | grep "host:" | head -1 | awk '{print $2}' | tr -d '"')
DB_PORT=$(grep -A 15 "mysql:" "$CONFIG_FILE" | grep "port:" | head -1 | awk '{print $2}' | tr -d '"')
DB_USER=$(grep -A 15 "mysql:" "$CONFIG_FILE" | grep "username:" | awk '{print $2}' | tr -d '"')
DB_PASS=$(grep -A 15 "mysql:" "$CONFIG_FILE" | grep "password:" | head -1 | awk '{print $2}' | tr -d '"')
DB_NAME=$(grep -A 15 "mysql:" "$CONFIG_FILE" | grep "database:" | head -1 | awk '{print $2}' | tr -d '"')

echo "正在初始化数据字典..."
echo "数据库: $DB_NAME"

# 执行 SQL
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" <<'EOF'

-- 清空现有数据（可选，根据需要决定是否保留）
-- DELETE FROM orbia_dictionary_item;
-- DELETE FROM orbia_dictionary;

-- ============================================
-- 1. 国家字典（Country）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('COUNTRY', '国家地区', '全球国家和地区列表，支持二级省份/州选择', 1)
ON DUPLICATE KEY UPDATE name='国家地区', description='全球国家和地区列表，支持二级省份/州选择';

SET @country_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'COUNTRY');

-- 美国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'US', 'United States', 1, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='United States';

SET @us_id = LAST_INSERT_ID();

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@country_dict_id, @us_id, 'US_CA', 'California', 1, 2, CONCAT(@country_dict_id, '/', @us_id), 1),
(@country_dict_id, @us_id, 'US_NY', 'New York', 2, 2, CONCAT(@country_dict_id, '/', @us_id), 1),
(@country_dict_id, @us_id, 'US_TX', 'Texas', 3, 2, CONCAT(@country_dict_id, '/', @us_id), 1),
(@country_dict_id, @us_id, 'US_FL', 'Florida', 4, 2, CONCAT(@country_dict_id, '/', @us_id), 1),
(@country_dict_id, @us_id, 'US_IL', 'Illinois', 5, 2, CONCAT(@country_dict_id, '/', @us_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 中国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'CN', 'China', 2, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='China';

SET @cn_id = LAST_INSERT_ID();

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@country_dict_id, @cn_id, 'CN_BJ', 'Beijing', 1, 2, CONCAT(@country_dict_id, '/', @cn_id), 1),
(@country_dict_id, @cn_id, 'CN_SH', 'Shanghai', 2, 2, CONCAT(@country_dict_id, '/', @cn_id), 1),
(@country_dict_id, @cn_id, 'CN_GD', 'Guangdong', 3, 2, CONCAT(@country_dict_id, '/', @cn_id), 1),
(@country_dict_id, @cn_id, 'CN_ZJ', 'Zhejiang', 4, 2, CONCAT(@country_dict_id, '/', @cn_id), 1),
(@country_dict_id, @cn_id, 'CN_JS', 'Jiangsu', 5, 2, CONCAT(@country_dict_id, '/', @cn_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 日本
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'JP', 'Japan', 3, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='Japan';

SET @jp_id = LAST_INSERT_ID();

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@country_dict_id, @jp_id, 'JP_13', 'Tokyo', 1, 2, CONCAT(@country_dict_id, '/', @jp_id), 1),
(@country_dict_id, @jp_id, 'JP_27', 'Osaka', 2, 2, CONCAT(@country_dict_id, '/', @jp_id), 1),
(@country_dict_id, @jp_id, 'JP_14', 'Kanagawa', 3, 2, CONCAT(@country_dict_id, '/', @jp_id), 1),
(@country_dict_id, @jp_id, 'JP_23', 'Aichi', 4, 2, CONCAT(@country_dict_id, '/', @jp_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 英国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'GB', 'United Kingdom', 4, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='United Kingdom';

SET @gb_id = LAST_INSERT_ID();

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@country_dict_id, @gb_id, 'GB_ENG', 'England', 1, 2, CONCAT(@country_dict_id, '/', @gb_id), 1),
(@country_dict_id, @gb_id, 'GB_SCT', 'Scotland', 2, 2, CONCAT(@country_dict_id, '/', @gb_id), 1),
(@country_dict_id, @gb_id, 'GB_WLS', 'Wales', 3, 2, CONCAT(@country_dict_id, '/', @gb_id), 1),
(@country_dict_id, @gb_id, 'GB_NIR', 'Northern Ireland', 4, 2, CONCAT(@country_dict_id, '/', @gb_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 德国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'DE', 'Germany', 5, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='Germany';

-- 法国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'FR', 'France', 6, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='France';

-- 加拿大
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'CA', 'Canada', 7, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='Canada';

-- 澳大利亚
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'AU', 'Australia', 8, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='Australia';

-- 韩国
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'KR', 'South Korea', 9, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='South Korea';

-- 新加坡
INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status)
VALUES (@country_dict_id, 0, 'SG', 'Singapore', 10, 1, CONCAT(@country_dict_id), 1)
ON DUPLICATE KEY UPDATE name='Singapore';

-- ============================================
-- 2. 性别字典（Gender）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('GENDER', '性别', '用户性别选项', 1)
ON DUPLICATE KEY UPDATE name='性别', description='用户性别选项';

SET @gender_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'GENDER');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@gender_dict_id, 0, 'ALL', 'All', 1, 1, CONCAT(@gender_dict_id), 1),
(@gender_dict_id, 0, 'MALE', 'Male', 2, 1, CONCAT(@gender_dict_id), 1),
(@gender_dict_id, 0, 'FEMALE', 'Female', 3, 1, CONCAT(@gender_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 3. 年龄段字典（Age Range）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('AGE_RANGE', '年龄段', '用户年龄段划分', 1)
ON DUPLICATE KEY UPDATE name='年龄段', description='用户年龄段划分';

SET @age_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'AGE_RANGE');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@age_dict_id, 0, 'AGE_13_17', '13-17', 1, 1, CONCAT(@age_dict_id), 1),
(@age_dict_id, 0, 'AGE_18_24', '18-24', 2, 1, CONCAT(@age_dict_id), 1),
(@age_dict_id, 0, 'AGE_25_34', '25-34', 3, 1, CONCAT(@age_dict_id), 1),
(@age_dict_id, 0, 'AGE_35_44', '35-44', 4, 1, CONCAT(@age_dict_id), 1),
(@age_dict_id, 0, 'AGE_45_54', '45-54', 5, 1, CONCAT(@age_dict_id), 1),
(@age_dict_id, 0, 'AGE_55_PLUS', '55+', 6, 1, CONCAT(@age_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 4. 语言字典（Languages）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('LANGUAGE', '语言', '支持的语言列表', 1)
ON DUPLICATE KEY UPDATE name='语言', description='支持的语言列表';

SET @lang_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'LANGUAGE');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@lang_dict_id, 0, 'EN', 'English', 1, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'ZH', 'Chinese', 2, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'ES', 'Spanish', 3, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'FR', 'French', 4, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'DE', 'German', 5, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'JA', 'Japanese', 6, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'KO', 'Korean', 7, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'PT', 'Portuguese', 8, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'RU', 'Russian', 9, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'AR', 'Arabic', 10, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'HI', 'Hindi', 11, 1, CONCAT(@lang_dict_id), 1),
(@lang_dict_id, 0, 'IT', 'Italian', 12, 1, CONCAT(@lang_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 5. 消费能力字典（Spending Power）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('SPENDING_POWER', '消费能力', '用户消费能力等级', 1)
ON DUPLICATE KEY UPDATE name='消费能力', description='用户消费能力等级';

SET @spending_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'SPENDING_POWER');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@spending_dict_id, 0, 'STANDARD', 'Standard', 1, 1, CONCAT(@spending_dict_id), 1),
(@spending_dict_id, 0, 'HIGH', 'High Spending Power', 2, 1, CONCAT(@spending_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 6. 操作系统字典（Operating System）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('OS', '操作系统', '设备操作系统类型', 1)
ON DUPLICATE KEY UPDATE name='操作系统', description='设备操作系统类型';

SET @os_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'OS');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@os_dict_id, 0, 'IOS', 'iOS', 1, 1, CONCAT(@os_dict_id), 1),
(@os_dict_id, 0, 'ANDROID', 'Android', 2, 1, CONCAT(@os_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 7. Android 系统版本字典（Android OS Versions）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('ANDROID_VERSION', 'Android版本', 'Android操作系统版本列表', 1)
ON DUPLICATE KEY UPDATE name='Android版本', description='Android操作系统版本列表';

SET @android_ver_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'ANDROID_VERSION');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@android_ver_dict_id, 0, 'ANDROID_14', 'Android 14', 1, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_13', 'Android 13', 2, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_12', 'Android 12', 3, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_11', 'Android 11', 4, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_10', 'Android 10', 5, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_9', 'Android 9 (Pie)', 6, 1, CONCAT(@android_ver_dict_id), 1),
(@android_ver_dict_id, 0, 'ANDROID_8', 'Android 8 (Oreo)', 7, 1, CONCAT(@android_ver_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 8. iOS 系统版本字典（iOS OS Versions）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('IOS_VERSION', 'iOS版本', 'iOS操作系统版本列表', 1)
ON DUPLICATE KEY UPDATE name='iOS版本', description='iOS操作系统版本列表';

SET @ios_ver_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'IOS_VERSION');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@ios_ver_dict_id, 0, 'IOS_17', 'iOS 17', 1, 1, CONCAT(@ios_ver_dict_id), 1),
(@ios_ver_dict_id, 0, 'IOS_16', 'iOS 16', 2, 1, CONCAT(@ios_ver_dict_id), 1),
(@ios_ver_dict_id, 0, 'IOS_15', 'iOS 15', 3, 1, CONCAT(@ios_ver_dict_id), 1),
(@ios_ver_dict_id, 0, 'IOS_14', 'iOS 14', 4, 1, CONCAT(@ios_ver_dict_id), 1),
(@ios_ver_dict_id, 0, 'IOS_13', 'iOS 13', 5, 1, CONCAT(@ios_ver_dict_id), 1),
(@ios_ver_dict_id, 0, 'IOS_12', 'iOS 12', 6, 1, CONCAT(@ios_ver_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 9. 设备品牌字典（Device Model/Brand）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('DEVICE_BRAND', '设备品牌', '设备制造商品牌列表', 1)
ON DUPLICATE KEY UPDATE name='设备品牌', description='设备制造商品牌列表';

SET @device_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'DEVICE_BRAND');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@device_dict_id, 0, 'APPLE', 'Apple', 1, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'SAMSUNG', 'Samsung', 2, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'GOOGLE', 'Google', 3, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'HUAWEI', 'Huawei', 4, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'XIAOMI', 'Xiaomi', 5, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'OPPO', 'OPPO', 6, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'VIVO', 'vivo', 7, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'ONEPLUS', 'OnePlus', 8, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'MOTOROLA', 'Motorola', 9, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'SONY', 'Sony', 10, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'LG', 'LG', 11, 1, CONCAT(@device_dict_id), 1),
(@device_dict_id, 0, 'NOKIA', 'Nokia', 12, 1, CONCAT(@device_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 10. 网络情况字典（Connection Type）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('CONNECTION_TYPE', '网络类型', '网络连接类型', 1)
ON DUPLICATE KEY UPDATE name='网络类型', description='网络连接类型';

SET @conn_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'CONNECTION_TYPE');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@conn_dict_id, 0, 'WIFI', 'WiFi', 1, 1, CONCAT(@conn_dict_id), 1),
(@conn_dict_id, 0, '2G', '2G', 2, 1, CONCAT(@conn_dict_id), 1),
(@conn_dict_id, 0, '3G', '3G', 3, 1, CONCAT(@conn_dict_id), 1),
(@conn_dict_id, 0, '4G', '4G', 4, 1, CONCAT(@conn_dict_id), 1),
(@conn_dict_id, 0, '5G', '5G', 5, 1, CONCAT(@conn_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 11. 时区字典（Time Zone）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('TIMEZONE', '时区', '全球时区列表', 1)
ON DUPLICATE KEY UPDATE name='时区', description='全球时区列表';

SET @tz_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'TIMEZONE');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, sort_order, level, path, status) VALUES
(@tz_dict_id, 0, 'UTC-12', 'UTC-12:00 (Baker Island)', 1, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-11', 'UTC-11:00 (American Samoa)', 2, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-10', 'UTC-10:00 (Hawaii)', 3, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-9', 'UTC-09:00 (Alaska)', 4, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-8', 'UTC-08:00 (Pacific Time - Los Angeles)', 5, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-7', 'UTC-07:00 (Mountain Time - Denver)', 6, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-6', 'UTC-06:00 (Central Time - Chicago)', 7, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-5', 'UTC-05:00 (Eastern Time - New York)', 8, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-4', 'UTC-04:00 (Atlantic Time)', 9, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-3', 'UTC-03:00 (Buenos Aires)', 10, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-2', 'UTC-02:00 (South Georgia)', 11, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC-1', 'UTC-01:00 (Azores)', 12, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+0', 'UTC+00:00 (London, GMT)', 13, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+1', 'UTC+01:00 (Paris, Berlin)', 14, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+2', 'UTC+02:00 (Cairo, Athens)', 15, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+3', 'UTC+03:00 (Moscow, Istanbul)', 16, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+4', 'UTC+04:00 (Dubai)', 17, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+5', 'UTC+05:00 (Pakistan)', 18, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+5_30', 'UTC+05:30 (India)', 19, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+6', 'UTC+06:00 (Bangladesh)', 20, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+7', 'UTC+07:00 (Bangkok, Jakarta)', 21, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+8', 'UTC+08:00 (Beijing, Singapore)', 22, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+9', 'UTC+09:00 (Tokyo, Seoul)', 23, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+10', 'UTC+10:00 (Sydney)', 24, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+11', 'UTC+11:00 (Solomon Islands)', 25, 1, CONCAT(@tz_dict_id), 1),
(@tz_dict_id, 0, 'UTC+12', 'UTC+12:00 (Fiji, New Zealand)', 26, 1, CONCAT(@tz_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 12. KOL报价方案类别（KOL Plan Type）
-- ============================================
INSERT INTO orbia_dictionary (code, name, description, status) 
VALUES ('KOL_PLAN_TYPE', 'KOL报价方案类别', 'KOL报价方案类别选项：基础版、标准版、高级版', 1)
ON DUPLICATE KEY UPDATE name='KOL报价方案类别', description='KOL报价方案类别选项：基础版、标准版、高级版';

SET @kol_plan_type_dict_id = (SELECT id FROM orbia_dictionary WHERE code = 'KOL_PLAN_TYPE');

INSERT INTO orbia_dictionary_item (dictionary_id, parent_id, code, name, description, sort_order, level, path, status) VALUES
(@kol_plan_type_dict_id, 0, 'basic', 'Basic', '基础版方案', 1, 1, CONCAT(@kol_plan_type_dict_id), 1),
(@kol_plan_type_dict_id, 0, 'standard', 'Standard', '标准版方案', 2, 1, CONCAT(@kol_plan_type_dict_id), 1),
(@kol_plan_type_dict_id, 0, 'premium', 'Premium', '高级版方案', 3, 1, CONCAT(@kol_plan_type_dict_id), 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description);

EOF

if [ $? -eq 0 ]; then
    echo "✓ 数据字典初始化成功！"
else
    echo "✗ 数据字典初始化失败！"
    exit 1
fi

