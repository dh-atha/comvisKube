<?php
header('Content-Type: application/json');
// add time sleep 1s to simulate processing time
sleep(5);
echo json_encode([
    "status" => "success",
    "message" => "Aplikasi Cloud-Native Berjalan Terbuka!",
    "environment" => "Kubernetes Lokal",
    "hostname" => gethostname(),
    "timestamp" => date('Y-m-d H:i:s')
]);
?>