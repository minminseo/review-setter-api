terraform {
    required_providers {
        cloudflare = {
            source = "cloudflare/cloudflare"
            version = "~> 5"
        }
        render = {
            source = "render-oss/render"
        }
    }

    # Terraform CloudでStateファイル保存
    cloud {
        organization = "minmin-dev"

        workspaces {
            tags = ["review-setter-api"]
        }
    }
}

# ============================

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

provider "render" {
  api_key  = var.render_api_key
  owner_id = var.render_owner_id
}

# ============================

# Cloudflare DNS Records
# Renderのデフォルトドメインへルーティング
resource "cloudflare_dns_record" "api_cname" {
  zone_id = var.cloudflare_zone_id
  name    = var.api_cname_name
  type    = "CNAME"
  content = var.api_cname_value
  ttl     = 1
  proxied = false
}

# ルートドメイン用のAレコード
resource "cloudflare_dns_record" "root_a" {
  zone_id = var.cloudflare_zone_id
  name    = var.root_a_name
  type    = "A"
  content = var.root_a_value
  ttl     = 1
  proxied = true
}

# Vercelのドメイン所有権確認用TXTレコード
resource "cloudflare_dns_record" "vercel_txt" {
    zone_id = var.cloudflare_zone_id
    name    = var.vercel_txt_name
    type    = "TXT"
    content = var.vercel_txt_value
    ttl     = 1
    proxied = false
}

# ============================

# Render Web Service
resource "render_web_service" "backend" {
  name           = var.render_service_name
  plan           = "free"
  region         = "singapore"

  # 起動コマンド: マイグレーション実行後にバッチ処理とAPIサーバーをバックグラウンドで同時起動
  start_command  = "./migrate -path ./migrations -database \"$DATABASE_URL\" -verbose up && ./batch_app & ./app"

  custom_domains = [
    { name = var.custom_domain_name }
  ]

  runtime_source = {
    native_runtime = {
      auto_deploy   = true
      runtime       = "go"
      repo_url      = "https://github.com/minminseo/review-setter-api"
      branch        = var.github_branch

      # golang-migrateのダウンロード、API/Batchバイナリのビルド
      build_command = "curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && chmod +x migrate && go build -tags netgo -ldflags '-s -w' -o app ./cmd/api && go build -tags netgo -ldflags '-s -w' -o batch_app ./cmd/batch"
    }
  }

# 環境変数
  env_vars = {
    "API_DOMAIN"        = { value = var.api_domain }
    "DATABASE_URL"      = { value = var.database_url }
    "ENCRYPTION_KEY"    = { value = var.encryption_key }
    "FE_URL"            = { value = var.fe_url }
    "GO_ENV"            = { value = var.go_env }
    "HMAC_SECRET_KEY"   = { value = var.hmac_secret_key }
    "PORT"              = { value = "8080" }
    "POSTGRES_DB"       = { value = var.postgres_db }
    "POSTGRES_HOST"     = { value = var.postgres_host }
    "POSTGRES_PORT"     = { value = "5432" }
    "POSTGRES_PW"       = { value = var.postgres_pw }
    "POSTGRES_USER"     = { value = var.postgres_user }
    "RESEND_API_KEY"    = { value = var.resend_api_key }
    "RESEND_FROM_EMAIL" = { value = "Review Setter <verify@minmindev.com>" }
    "SECRET"            = { value = var.secret }
  }
}