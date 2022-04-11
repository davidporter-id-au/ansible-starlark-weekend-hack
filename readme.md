## weekend hack - ansible with starlark preprocessing

This is a quick hack to play with starlark and see how compatible it might be
to use with ansible. This is a weekend-hack and not production code.

### Problem

Ansible is fine as a task runner, and has a fairly powerful... I don't want to
call it 'inheritance' model, but something akin to that which allows for strong
naming consistency in large infra projects with ease.

But there always comes a time when you end up fighting yaml, which is both
inflexible to a point and sometimes fairly insane with the complexity possible
(particularly in ansible where you're doing preprocessing with the templating
engine anyway)

Moreover, there's a whole bunch of extremely annoying problems with ansible's
looping which continue to baffle-me, but are probably just legacy holdover. 

### Idea:

Not mine, it's Mantas', though I've quite possibly bastardized it: Could we use
starlark as a preprocessor (yes, I know ansible already has one)? If so, what
would it look like?

In practice, this means: 

1. Using Starlark based config as the source of truth
2. Generating yaml output with a questionable module system thats... kind of
   specified in starlark (so many questionable ideas here). 

   The idea being that each starlark file has an entrypoint function called `module` whose purpose
   is to be a pure function, which just returns data. The data returned becomes
   whatever it's ansible equivalent is. For `group_vars` it probably should
   just be some maps of config. For tasks in the `roles` folder it might be
   ansible blocks and tasks (with sane looping?)

### Quick-start: 

The idea, as mentioned, is that starlark's the source of truth and the yaml files are just generated artifacts sitting alongside mostly (although, there's nothing stopping you from using either). 

To run ansible-playbook, execute the go app. Eg: `go run astar.go example-playbook.yaml`

This will:

1. generate the ansible yaml
2. execute ansible-playbook, passing whatever args you give it

### Credit

most of the actual `starlark` -> `yaml` converstion was [written by `cirrus-ci`](https://github.com/cirruslabs/cirrus-cli/pull/46/files#diff-fd961f8f67870410b5925d977385825c70a2da811309b101a719b383ee2d8a04) and I've just simplified it.

