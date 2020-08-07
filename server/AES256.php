<?php
function encrypt($text, $passphrase)
{
    $salt = openssl_random_pseudo_bytes(8);

    $salted = $dx = '';
    while (strlen($salted) < 48) {
        $dx = md5($dx . $passphrase . $salt, true);
        $salted .= $dx;
    }

    $key = substr($salted, 0, 32);
    $iv = substr($salted, 32, 16);

    return base64_encode('Salted__' . $salt . openssl_encrypt($text . '', 'aes-256-cbc', $key, OPENSSL_RAW_DATA, $iv));
}

function decrypt($encrypted, $passphrase)
{
    $encrypted = base64_decode($encrypted);
    $salted = substr($encrypted, 0, 8) == 'Salted__';

    if (!$salted) {
        return null;
    }

    $salt = substr($encrypted, 8, 8);
    $encrypted = substr($encrypted, 16);

    $salted = $dx = '';
    while (strlen($salted) < 48) {
        $dx = md5($dx . $passphrase . $salt, true);
        $salted .= $dx;
    }

    $key = substr($salted, 0, 32);
    $iv = substr($salted, 32, 16);

    return openssl_decrypt($encrypted, 'aes-256-cbc', $key, true, $iv);
}