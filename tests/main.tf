resource null_resource date {
  provisioner local-exec {
    command = "date"
  }
}