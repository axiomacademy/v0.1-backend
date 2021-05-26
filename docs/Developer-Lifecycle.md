# Developer Lifecycle :recycle: 
This section tells you how to contribute to the codebase in a productive way, using all the tools GitLab provides for you to use. It details the formal processes involved in creating an update to the code without accidentally messing up other people's code, and it also tells you how to get your code merged.

## Protected branches and `master`
The default git branch is called `master` :star:. By default for security and code compliance reason, the `master` branch is a protected branch. That means that only people with *Maintainer* authority and above can push directly to the master branch. You also need their approval to merge any branches into `master` as you likely would be doing in the future

## Step 1: Adopt an issue :exclamation: 
When you get started, it is likely that someone will assign you an issue. An issue has a name for example `feat: creating chat service` and a number #3. Read the issue carefully and understand the problem before you start. If in doubt, clarify with the person who assigned you the issue :smile: 

## Step 2: Create a branch :tanabata_tree: 
When you are working on your code, you want to avoid accidentally trampling on other people working on a similar portion of the code. The way we manage these conflicts in git is through branches. I won't bother going into too much detail you can find more information on how to create and manage branches [here](https://git-scm.com/book/en/v2/Git-Branching-Basic-Branching-and-Merging)

When you are working with branches, one of the most important things is naming. We follow semantic branch naming conventions. So if the issue you are working on is `feat: creating chat service` your branch should be appropriately named `feat/creating-chat-service`. Spaces are illegal in branch names so we use the / and - symbols to compensate. 

## Step 3: Create your commits :clap: 
Now you are ready to actually start coding. As you write your code, pick logical intervals to pause and commit your code to version control. This not only allows you to reconsolidate your thinking, but it makes it easier to revert and see your progress in case you make any mistakes. It also simplifies later debugging as an error might be able to be isolated to a single commit. 

Naming your commits semantically is as important as naming your branches. We follow a commit naming format as follows `commit_type: commit message here`. The valid commit types are:

```
feat -> represents features
docs -> represents documentation only commits
refactor -> represents refactoring work
merge -> represents merges from other branches
```

Therefore, an example commit message might look like `feat: implementing chat service`

## Step 4: Create your merge request 😲
Finally, you are done writing your code and ready to merge it. Head on over to the merge request, and request for a merge of your branch into `master`. That would mean that `master` is the destination branch and your branch is the source. 

Name your merge request semantically too, **follow the naming convention used for your issues.** Just as cool is putting a link to your issue in the merge commit. This allows code reviewers to easily track what issues you are solving. That can be done with a # symbol, so, for example, closes #3 is a valid link that would automatically close issue 3, should the merge request be merged.

That's it! Those are the basic rules and have fun developing :smiling_imp: 
