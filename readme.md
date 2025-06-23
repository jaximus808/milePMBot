# MilestonePM Bot

![alt text](https://github.com/jaximus808/milePMBot/blob/main/assets/MilestonePM.png "Random Logo LOL")

Thanks for checking out the MilestonePM bot! 😊

This Discord bot helps manage projects, milestones, and tasks within your Discord server through text-commands. It provides a comprehensive project management system with role-based permissions and milestone tracking. 

With text-commands, updating and tracking progress is much quicker compared to drag and drop. No more navigating confusing UI, just run simple commands that are easier to navigate thanks to Discords API. 

## Concepts and Usuage

MilestonePM works by mapping project to discord server categories and its channels. This means that any channel within an active project's category can run commands to interact with that project

![alt text](https://github.com/jaximus808/milePMBot/blob/main/assets/milestonePMReference.png "Visual Category Pic")


Within a project, there are 2 main terms

### Milestones

These tie to big "strides" within a project, and are used to track the general direction of a project. 
Each project will have an active milestone, meaning all tasks and work will be tied to that milesstone. Commands can be used to move between milestones directional based on its due date.

In Jira Terms these are essentially Epics

### Tasks

Tasks are the specific work that a member is doing towards completing an active milestone

Taks go through 4 phases

1. Creation - Tasks are created with its description and expectation, but are not assigned yet (Backlog)
2. In Progress - When a task is assigned it is deemed in progress, the assigned can update the project with progress reports
3. In Review - When the assigned user marks a project as complete, the bot will notify the assigner to review the work
4. Done - Once a project is approved by the assigner, the task is deemed as complete and is closed 

In Jira terms these are essentially stories

## Permission Levels

Following permissions are outlined below in order of most -> least privilege. 

- **OWNER**: Can end, resume, and move projects
- **ADMIN**: Can manage milestones, approve/reject any tasks
- **LEADS**: Can create and assign tasks, review tasks they assigned

**DISCORD SERVER ADMIN**: This is not really a milestonePM role, but a discord permission level required to start project within a server

Any commands that don't have a required permissions outlined is able to be used by anyone

A + indicates that role or higher in privilege can access that command

### Usuage 

The workflow to using MilestonePM is through discord-slash commands. This allows users to interact with a project, update progress, and track work. 

In upcoming updates, a web UI view will be avilable to allow a kaban style board for tasks. This webpage will also be the home center for other general usuages for the app. More info coming soon

The following documentation convers specifics for using the bot. You may also use /help command to get a concise version of these documents

## Command Documentation Format

**Arguments**: `[argname]` - Command arguments where the brackets indicate to put a value and `argname` is the name of the argument you are setting. Format: `argname:value`

**Example**: `/help [cmd]` - The help command expects an argument "cmd" with the command name  
Valid usage: `/help cmd:project` - Gives more info about project commands

## Autocomplete

- Please use autocomplete! This will give you the available options for a command and speed up writing commands
- We'll continue working to add more autocomplete in the future

*Note*: Spaces are okay with command arguments thanks to Discord's API 😊

## General Commands

### `/help`
Prints the main helper message

### `/help [command]`
Gives more detailed information for a specific sub-command

**Example**: `/help cmd:project` - gives information about project commands

## Project Commands

The `/project` command handles project settings and control with different permission levels.

### Project Control (DISCORD SERVER ADMIN)

#### `/project start [msname] [msdate] [msdesc]`
- Starts a project in the given Discord channel category
- Initiates the project's first milestone with the given values
- The initial milestone becomes the active milestone
- Date format: MM/DD/YYYY
- Your account/role must be authorized to start a project

### Project Control (OWNER)

#### `/project end [projectref]`
- Ends a project and removes it from active projects
- Requires inputting the active projectref to confirm the project's ending

#### `/project resume [projectref]`
- Resumes a project that was ended and moves it to the current channel's category
- Your Discord account must have been an owner in that project

#### `/project move [projectref]`
- Moves an active project to another Discord category
- All tasks, milestones, and roles are maintained during the move
- Must be the owner of the project

### Project Settings (OWNER)

#### `/project set [setting] [value]`
Updates project settings. Available settings include:
- Changing output channel
- Update project description
- Toggle sprints
- Sprint message
- Sprint length
- Toggle sprint pings

### Project Roles (ADMIN+)

#### `/project role [op] [user] [role]`
Manages user roles within the project
- `role` parameter is optional for remove operations

### Project Info

#### `/project info`
Displays information for the current project

## Milestone Commands

The `/milestone` command handles creation and movement of milestones for the current project (ADMIN + LEADS ONLY).

### Milestone Control (ADMIN+)

#### `/milestone create [msname] [msdate] [msdesc]`
- Creates a milestone for the project with the given arguments
- Returns the milestone reference tied to this milestone
- Note: Does not make this milestone active

#### `/milestone move [direction]`
- Moves the active milestone to the previous or next milestone
- Direction options: `next` or `prev`
- Movement is determined by due date order
- Example: If current milestone is due 06/20/2025 and next is 07/25/2025, `/milestone move direction:next` moves to the 07/25/2025 milestone

#### `/milestone delete [taskref]`
- Removes the milestone and all associated tasks
- Requires confirmation message
- **This is permanent and cannot be undone!**

### Milestone Info

#### `/milestone map`
Creates a milestone map showing task status:
- Unassigned tasks
- In progress tasks
- In review tasks
- Completed tasks

## Task Commands

The `/task` command handles creating, assigning, progress tracking, completing, and approvals for tasks within a milestone.

### Starting Tasks (LEADS+)

#### `/task create [name] [desc]`
- Creates a new task for the current milestone
- Returns a task reference for use in other commands

#### `/task assign [user] [taskref] [expectation]`
- Assigns a task to a user
- Task expectations can use either AGILE story points or due dates

### Task Progress

#### `/task progress [taskref] [desc]`
- Creates a progress update for the user's assigned task
- Description outlines the current progress

#### `/task complete [taskref] [desc]`
- Marks a task as complete and ready for review
- Pings the assigner to review the task

### Reviewing Tasks (LEADS+)

#### `/task approve [taskref]`
- Marks a task as approved and done for the current milestone
- Leads can only approve tasks they assigned to normal members
- Admins can approve any tasks, including tasks of other admins

#### `/task reject [taskref] [desc]`
- Marks a task as not approved and returns it to in-progress status
- Description provides feedback on what needs to be fixed
- Leads can only reject tasks they assigned to normal members
- Admins can reject any tasks, including tasks of other admins

### General Task Commands

#### `/task list`
Lists all tasks for the current milestone and displays their status

#### `/task list [user]`
Lists tasks assigned to a specific user, showing tasks that are:
- In progress
- In review
- Complete

## Current Planned Updates

### CRITICAL

- A Web version to allow a kaban style UI to visualize projects
- Deploying code onto VPS server for hosting, with CI/CD
- More verbose and cleaner discord messages, info (such as stats), and just code in general
- Writing actual tests

### Important

- Sharding
- More bots to other platforms, such as slack

## Contributing

Contributions are welcome! Due to how new the project is no active proccedure is made to support contributions, but pleae feel free to open any issues for desired features.

## Using the Bot

- This bot will soon be avialable for usuage with an invite link, as I will set up hosting to have the bot online 24/7
- I plan to make a free tier, and then a paid tier for more features and support
- If you wish to use the bot early or self host, feel free to clone this repo. However you will need to use your supabase database, or feel free to rewrite the helper methods to use your own database model. 


