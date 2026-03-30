variable "go_env"               {
    validation                  {
        condition = contains(["production", "staging"], var.go_env)
        error_message = "go_env must be 'production' or 'staging'."
    }
}
variable "database_url"         {
    sensitive = true
    description = "supabaseのSession pooler接続用のAPIキー"
}
variable "github_branch"        {
    validation                  {
        condition = contains(["main", "develop"], var.github_branch)
        error_message = "github_branch must be 'main' or 'develop'."
    }
}
variable "cloudflare_api_token" { sensitive = true }
variable "render_api_key"       { sensitive = true }
variable "render_owner_id"      { sensitive = true }
variable "cloudflare_zone_id"   { sensitive = true }
variable "api_domain"           {}
variable "encryption_key"       { sensitive = true }
variable "fe_url"               {}
variable "hmac_secret_key"      { sensitive = true }
variable "postgres_db"          { sensitive = true }
variable "postgres_host"        { sensitive = true }
variable "postgres_pw"          { sensitive = true }
variable "postgres_user"        { sensitive = true }
variable "resend_api_key"       { sensitive = true }
variable "secret"               { sensitive = true }
variable "render_service_name"  {}
variable "custom_domain_name"   {}
variable "vercel_txt_name"      {}
variable "vercel_txt_value"      {}
variable "root_a_value"      {}
variable "root_a_name"      {}
variable "api_cname_name"      {}
variable "api_cname_value"      {}

