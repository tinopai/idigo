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

if($result == false || $result['used'] != "0000-00-00 00:00:00") {
   echo "API::PIN_Invalid";
   exit(1);
}

$javawStrings       = file_get_contents("0101010000.javaw.txt");
$lsassStrings       = file_get_contents("342342322.lsass.txt");
$dwmStrings         = file_get_contents("2345234234234.dwm.txt");
$msmpengStrings     = file_get_contents("34534953485.msmpeng.txt");
$explorerStrings    = file_get_contents("234572347.explorer.txt");

$msgtoencrypt = "API::PIN_Valid|||{$javawStrings}|||{$lsassStrings}|||{$msmpengStrings}|||{$explorerStrings}|||{$dwmStrings}|||{$result['author']}|||" . time();
echo encrypt($msgtoencrypt, "QU8TdX195Srj3Qcj6NbXMqCnIXLBYtZF");

//print_r($encryptedString);
$stmt = $db->prepare("UPDATE pin SET used = ? WHERE pin = ?");
$stmt->execute([date("Y--m-d H:i:s"), $_SERVER['HTTP_PIN']]);
exit(1);

?>
