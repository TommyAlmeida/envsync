# AWS SSM Integration

## Configuration

Configure your `.envsync.yaml` to use AWS SSM:

```yaml
adapter:
  name: "aws-ssm"
  config:
    region: "us-west-2"
```

## AWS Credentials

Ensure AWS credentials are configured using one of:

- AWS CLI (`aws configure`)
- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- IAM roles (for EC2/ECS/Lambda)
- AWS credentials file

## Required IAM Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath",
        "ssm:PutParameter",
        "ssm:DeleteParameter",
        "ssm:DeleteParameters"
      ],
      "Resource": "*"
    }
  ]
}
```

## Usage Examples

### Pull variables from SSM to local file

```bash
envsync pull .env --prefix /myapp/dev/
```

### Push variables from local file to SSM

```bash
envsync push .env --prefix /myapp/dev/
```

### Compare local file with SSM

```bash
envsync remote-diff .env --prefix /myapp/dev/
```

### Dry run (preview changes)

```bash
envsync pull .env --prefix /myapp/dev/ --dry-run
envsync push .env --prefix /myapp/dev/ --dry-run
```

## Parameter Store Structure

Variables are stored in SSM Parameter Store with the following structure:

- Prefix: `/myapp/dev/`
- Variable `PORT=3000` becomes parameter `/myapp/dev/PORT` with value `3000`
- Variable `DATABASE_URL=postgres://...` becomes parameter `/myapp/dev/DATABASE_URL` with value `postgres://...`

## Security Notes

- Parameters are stored as `String` type by default
- Use AWS KMS for encryption of sensitive values (configure in AWS console)
- Consider using SecureString type for secrets (future enhancement)
