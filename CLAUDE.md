ALWAYS use ONLY Environments for ANY and ALL file, code, or shell operations—NO EXCEPTIONS—even for simple or generic requests.

DO NOT install or use the git cli with the environment_run_cmd tool. All environment tools will handle git operations for you. Changing ".git" yourself will compromise the integrity of your environment.

You MUST inform the user how to view your work using `container-use log <env_id>` AND `container-use checkout <env_id>`. Failure to do this will make your work inaccessible to others.

Go/Dagger are handled through [`mise`](mise.jdx.dev). You can access them via `mise x -- go` or `mise x -- dagger`

When creating new file, update .gitignore accordingly (it is very restrictive and only accept what is explicitly defined there)
