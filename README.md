# Arithland Telegram

## Local Development

**1.** `cd` somewhere you want to keep the project.
```shell
mkdir repos
cd repos
```

**2.** _(Only first time)_ clone the git repository.
```shell
git clone git@github.com:avistopia/arithland-telegram.git
```

**3.** _(Only first time)_ Copy `.env.template` to `.env`. Add secret values into the `.env` file. Make sure to never commit
`.env` file, because the project is open source and the source files hosted on GitHub are available to the public.

**4.** `cd` inside the repository directory.
```shell
cd arithland-telegram
```

**5.** Start the container.
```shell
docker compose -f docker-compose-local.yaml up -d
```

**6.** Exec into the container.
```shell
docker compose -f docker-compose-local.yaml exec arithland-telegram bash
```

**7.** Run the project inside the container.
```shell
make run
```

**8.** You can stop the bot using ctrl+C (cancel). Wait a few seconds for it to shut down gracefully.

**9.** You can exit the exec shell using `exit` command or ctrl+D. When you exit the exec shell, the container will
remain running.

**10.** Whenever you're done, you can stop the container.
```shell
docker compose -f docker-compose-local.yaml down -d
```

**Notes**

- You can start the container again using the command in step **5**.
- Whenever you make changes to the code, make sure to stop (**8**), and run the bot again (**7**), for your changes 
  to affect.
- Whenever you modify `.env` values, make sure to stop (**10**) and start the container again (**5**), for it to affect.

## Contribution

**1.** _(Only first time)_ setup GitHub SSH key. (Refer to the [GitHub SSH Key Doc](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account?platform=mac))

**2.** Make your changes. Run test and lint inside the docker container.
```shell
make lint
make test
```

**3.** Add and commit your changes. (Refer to the [GoLand Git Doc](https://www.jetbrains.com/help/go/commit-and-push-changes.html))

**4.** Make sure to pull before pushing your commit. Make sure to use rebase mode and not the default merge mode.
This can be done via the IDE GUI as well as the CLI.
```shell
git pull --rebase
```

**5.** Push your changes.

**6.** Monitor the triggered pipeline in GitHub actions. ([Arithland Bot GitHub Actions](https://github.com/avistopia/arithland-telegram/actions))

**7.** If the triggered pipeline fails, this means either tests or lints failed, or some error happened when building
and deploying the project. The exact error can be found in the pipeline logs.

**8.** If the triggered pipline finishes successfully, your changes will be available after a few minutes on the server.

**Notes**

- If receiving lint issues, running the following command inside the docker container can help fix some of the issues
  automatically.
```shell
make lint-fix
```