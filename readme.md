## weekend hack - ansible with starlark preprocessing

This is a quick hack to play with starlark and see how compatible it might be to use with ansible. This is a weekend-hack and not production code.

### Problem

Ansible is fine as a task runner, and has a fairly powerful... I don't want to call it 'inheritance' model, but something akin to that which allows for strong naming consistency in large infra projects with ease.

But there always comes a time when you end up fighting yaml, which is both inflexible to a point and sometimes fairly insane with the complexity possible
(particularly in ansible where you're doing preprocessing with the templating engine anyway)

Moreover, there's a whole bunch of extremely annoying ugliness with ansible's looping which continue to baffle me.

### Idea:

Not mine, it's Mantas', though I've quite possibly bastardized it: **Could we use starlark as (another) preprocessor**? If so, what would it look like?

In practice, this means: 

1. Using Starlark based config as the source of truth
2. Generating yaml output from starlark with pure functions
3. Running ansible on the generated yaml as before. 

**Modules and parsing code...?**

Starlark has a `load` keyword module system thats... kind of specified, but is left to the application / lanuage implementer to actually make happen. So, in order to make that work it's necessary to parse the starlark, and then provide some functionality on top of the language for things such as loading files. 

So far so good, but the next questions is how to get stuff out of the starlark file, since it's a fully turing-complete language, and contains far more than just data. 

Following on with the idea by the code from cirrus-ci that I was shamelessly repurposing, I came to the conclusion that it's probably best to have some kind of entrypoint into the starlark files in order to separate out the function configuration from the pure exported data. So I just made up the convention of using relying on a function called `module` in whatever starlark file's there, and using it's exports.

The data returned becomes whatever it's ansible equivalent is. For `group_vars` it probably should just be some maps of config. For tasks in the `roles` folder it might be ansible blocks and tasks (with sane looping?)

### Quick-start: 

The idea, as mentioned, is that starlark's the source of truth and the yaml files are just generated artifacts sitting alongside mostly (although, there's nothing stopping you from using either). 

To run ansible-playbook, execute the go app. Eg: `go run astar.go example-playbook.yaml`

This will:

1. generate the ansible yaml
2. execute ansible-playbook, passing whatever args you give it

### Credit

most of the actual `starlark` -> `yaml` converstion was [written by `cirrus-ci`](https://github.com/cirruslabs/cirrus-cli/pull/46/files#diff-fd961f8f67870410b5925d977385825c70a2da811309b101a719b383ee2d8a04) and I've just simplified it.

