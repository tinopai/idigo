<?php

date_default_timezone_set('UTC');
$version = "v0.1a";

if($_SERVER['HTTP_USER_AGENT'] != "idigo ({$version}): client") {
    echo "Error";
    exit(1);
}

require('config.php');
require('AES256.php');

$_SERVER['HTTP_ST'] = Decrypt($_SERVER['HTTP_ST'], "4Oe7EmckEVKuogjcoLQWAaVhJAkIs6PT");
if($_SERVER['HTTP_ST'] != "6mYiPz5u1WO62MM80kyZpBEQXHlD0Ho8nEwqzcegSj1ZeqqeplZHcRAry4e0lKWyUnn7VKiSOdrWR897BGtxaxLJuAJze96mza2") {
    echo "Error";
    exit(1);
}

$_SERVER['HTTP_PIN'] = Decrypt($_SERVER['HTTP_PIN'], "PINMY8n7aVQ7WX03aAHV8mbzFBBsEWIO");
if(!isset($_SERVER['HTTP_PIN']) || strlen($_SERVER['HTTP_PIN']) != 6 || preg_match("/[^0-9]/i", $_SERVER['HTTP_PIN']) == 1 ) {
    echo "Error";
    exit(1);
}
if(!isset($_SERVER['HTTP_RESULTS'])) {
    echo "Error";
    exit(1);
}

function generateRandomString($length = 10) {
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $charactersLength = strlen($characters);
    $randomString = '';
    for ($i = 0; $i < $length; $i++) {
        $randomString .= $characters[rand(0, $charactersLength - 1)];
    }
    return $randomString;
}

$stmt = $db->prepare("SELECT * FROM pin WHERE pin = ?");
$stmt->execute([$_SERVER['HTTP_PIN']]);
$result = $stmt->fetch(PDO::FETCH_ASSOC);


$decryptedString = explode("|||", decrypt($_SERVER['HTTP_RESULTS'], "Uv21t0oNjUHlHSiMcTO8cIsWYa77vo7y"));
if(substr($decryptedString[1], 0, strlen("IDIGO")) != "IDIGO") {
    echo "Error";
    exit(1);
}

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL,"http://127.0.0.1:51335");
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_POSTFIELDS,$vars);  //Post Fields
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

$headers = [
    "Cheats: {$decryptedString[2]}",
    "Username: {$decryptedString[3]}",
    "Channel: {$result['channel']}",
    "Author: {$result['author']}",
    "Recycle: {$decryptedString[4]}",
    "Build: {$decryptedString[5]}"
];

curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);

$server_output = curl_exec ($ch);

curl_close ($ch);

echo Encrypt("API::VALID_REQUEST|||IDIGO|||" . time() . "|||" . generateRandomString(64) . "|||" . $result['channel'] . "|||" . $result['author'], "Uv21t0oNjUHlHSiMcTO8cIsWYa77vo7y");
$stmt = $db->prepare("DELETE FROM pin WHERE pin = ?");
$stmt->execute([$_SERVER['HTTP_PIN']]);
exit(1);

?>