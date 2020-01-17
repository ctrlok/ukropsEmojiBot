// Install dependencies and copy function to the folder with dependencies
resource "null_resource" "build_slackConnector" {
  triggers = {
    scriptsha     = join(",", [ for k in fileset(path.module, "slackConnector/*.go") : join(":", [k, filebase64sha256(k)]) ])
  }
  provisioner "local-exec" {
    command = "cd slackConnector; GOOS=linux GOARCH=amd64 go build -o slackConnectorBin ."
  }
}

data "null_data_source" "wait_for_build_slackConnector" {
  inputs = {
    # This ensures that this data resource will not be evaluated until
    # after the null_resource has been created.
    lambda_exporter_id = null_resource.build_slackConnector.id

    # This value gives us something to implicitly depend on
    # in the archive_file below.
    source_file = "${path.module}/slackConnector/slackConnectorBin"
  }
}

data "archive_file" "lambda" {
  output_path = "slackConnector.zip"
  type        = "zip"
  source_file = data.null_data_source.wait_for_build_slackConnector.outputs["source_file"]
}