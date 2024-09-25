// Define an environment named "local"
env "prod" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "file://store/mysql/schema.sql"

  // Define the URL of the database which is managed
  // in this environment.
  url = "mysql://fin:1234@localhost/fin"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "docker://mysql/8/fin"
}