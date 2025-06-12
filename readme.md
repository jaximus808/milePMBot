# MilePM Bot (Mile Project Management)

This project aims to allow for easier tracking of a project through a discord bot

Instead of slowly clicking and dragging tickets, just type commands and spin them up!

Much faster :)

# Some Terms and concepts

## Projects

Projects contain milestones that encompases the start to end of a project. The project also contains all data related to tasks, milestones, settings etc. 

## Milestones

Milestones are big top level goals within a project, each project will be on a certain milestone. 
Every task created and interacted on will be assigned to the project's current milestone. Admins dictate which milestone a project is on. 

## Tasks

Tasks are assigned to users and are given a due date, along with a description. 

When a person finishes a task its then sent for review by a upper role, admins dont need thier work review when finished

Tasks are tied to milestones, and are tracked within said milestone

## Sprints

By default two weeks dedicated to a quanitity of goals. All tasks will be assigned to some sprint and a message of TODOs for a specific sprint will be posted by a day offset

At the end of a sprint, the program will mark what was accomplished, and if a new milestone was met, etc

## Project Members

This is the default role for users within a project. These users can only update progress on a task, mark that they are finished, and print any data about a project. 

## Leads 

These users can access default commands, and are also able to assign tasks and approve tasks for a given project member. 

## Admins 

These users and full control of a project, and are able to assign tasks to leads. These users can also modify settings of a project and determine which milestone a project is currently on

# Usage 

## Set Up (DONE: NEEDS TESTING)

Each project will need to be set up through these following commands 

``` 
 !project start [initial milestone name] [initial milestone deadline] [inital milestone desc]
```
This will start the project with default settings, where sprint 0 starts the current day its ran (pushed back to 8am)

By default sprints are turned off

Only privilleged users can run this command, (and i'll figure out a way to set this, maybe some form of auth). By default the output channel will be the channel the command is ran 

## Settings

For project admins, you can adjust default settings here with !project set ...


### Setting Project Description Messages (WIP)

Set the project description! I recommend just link a google doc

``` 
 !project set project_desc [desc]
``` 
### Setting Custom Sprint Messages (WIP)

Change the sprint msg 

``` 
 !project set sprint_msg [mesg here no brackets]
``` 

### Toggle Project Sprint Pings (WIP)

Enable or disable @here pings for sprint messages 

``` 
 !project set sprint_ping [off | on]
``` 

### Toggle Project Automatic Sprint Messages (WIP)

Enable or disable end of sprint messages all together

``` 
 !project set sprint_auto [off | on]
``` 

### Change sprint week length (WIP)

Change the duration of a sprint (and when to auto remind of given sprints tasks)

``` 
 !project set sprint_length [#]
``` 

### Set Output Channel (WIP)

Set the channel for the bot to output sprint messages and other logs to the current id

``` 
 !project set output 
``` 

### Toggle Sprints (WIP)

If sprints is something you do not want for your project, then feel free to disable it! All tasks will only use dates, so anything commands using sprints will not work

``` 
 !project set sprints [off | on]
``` 

### Add Admin (WIP)

Add admin to a differnet user 

``` 
 !project add admin @user
``` 

### Remove Admin (WIP)

Remove admin from a user

``` 
 !project remove admin @user
``` 

For the sake that this isn't that deep, all admins have equal perms, meaning any admin can remove any other admin within a project. DONT BE DUMB PLEASE

### Add new Lead (WIP)

Add a new user to be a lead under a certain sub team

```
 !project add lead @user 
``` 

so only these users can set tasks for users under the given subteam 

### Remove Lead (WIP)

Removes user as project lead

``` 
 !project rm lead @user
``` 
## Project Control (ADMIN ONLY)

### Complete Project (WIP)

``` 
 !project complete
``` 

The bot will first check if all tasks in milestones are complete, and if not identify those milestones not done and ask if they wish to move foward, use !select [yes | no] to move on

The bot will then ask again if they are sure

use !select [yes | no] to make your decision

once done the project is marked as complete! A project report will be generated and you'll be able to view the projects progress and stats (and who did what)

### Pause Project (WIP)

``` 
 !project pause 
``` 

All work and auto stuff scheduled like sprints will be paused, no commands will work except for project control and project setting commands 

The command will prompt you to confirm action

### Resume Project (WIP)

``` 
 !project resume
``` 

Resumes a project, and updates you on any tasks that were missed due to the pause 

Command will also prompt confirmation

## Project Info (anyone)

Prints info out about a project set by settings

``` 
 !project info
``` 

## Project Team

prints out the team, roles, and stats

``` 
 !project team
``` 


## Milestones

### Create new milestone (Admin ONLY) (WIP)

Creates a new milestone and prints its id

``` 
 !milestone create [milestone name] [milestone deadline] [milestone desc]
``` 

### Delete new milestone (Admin ONLY) (WIP)

Delete a given milestone, CANNOT BE THE CURRENT MILESTONE

``` 
 !milestone delete [milestone id]
``` 

### Updating a milestone (ADMIN ONLY) (WIP)

Updating Deadline 

``` 
 !milestone update deadline [milestone id] [milestone deadline]
``` 

Updating desc 

``` 
 !milestone update desc [milestone id] [milestone desc]
``` 

Updating name 

``` 
 !milestone update name [milestone id] [milestone name]
``` 

### Move To Next Milestone (ADMIN ONLY) (WIP)

Moves to the next milestone if one exists 

``` 
 !milestone move next
``` 

Moves to the previous milestone if one exists

``` 
 !milestone move prev
``` 

Moves to a given milestone

``` 
 !milestone to [milesstone id]
``` 

### Milestone Info (anyone) (WIP)

Current milestone

``` 
 !milestone info
``` 

Given milestone

``` 
 !milestone info [milestone id]
``` 

### Overall Milestone Map (WIP)

Gets a map of the projects milestone map

``` 
 !milestone map
``` 

## Tasks (Steps)

For tasks commnads, for commands with leads and amdins they follow this permission structure

Member cannot modify any tasks 

Lead can only modify their tasks and their subordinate tasks

Admin can modify any tasks

Below will list out the usuage for tasks for a given project

### Create Tasks (Leads and Admins) (WIP)

Create a task for the current milestone

``` 
 !task create [task name] [desc] 
``` 

This command will give you a task_ref that has no spaces, where you will use to assign the task later

It will also give you task id, which you can use later

### Assign Tasks (Leads and Admins) (WIP)

Will assign a task for a user and will also give you a task id

A user cannot be given a task with the same name

Assigning a prev created task 

```
 !task assign @user [task_ref] [due_date]
```

Assign a task: 


Date version: 

``` 
 !task assign @user [task_ref] [due date] 
``` 

Sprint version (the task still has a date under the hood, and is due at the end of said sprint): 

``` 
 !task assign @user [task_ref] [story points] 
``` 

If a user is already assigned the task ref, and you assign a new user, the original user will no longer be assigned to that task

### Remove Task (Leads and Admins) (WIP)

Removes task from the project

``` 
 !task remove @task [task_id]
``` 

### Update deadline of Task (Leads and Admins) (WIP)

Update a task's deadline

date method: 

``` 
 !task update deadline [task_ref] [new date]
``` 


### Mark as Complete (complex perms) (WIP)

Approves a task after its marked as complete by its subords, and updates the task as complete

``` 
 !task approve [task_id]
``` 

### Dissaprove (complex perms) (WIP)

Dissaproves the tasks and makes sets it to not being done, same permissions as approve

``` 
 !task reject [task_id] [needed fixes]
``` 

### Retrieve Current Tasks (Any) (WIP)

Gets the user's current assigned tasks 

Get all tasks

``` 
 !task list
``` 

Get @user's tasks do

``` 
 !task list @user
``` 

### Update Progress on Task (Any) (WIP)

Update the current progress of a task

``` 
 !task progress [task_ref] [desc]
``` 




### Set to done on Task (Any) (WIP)

Update the current progress of a task

``` 
 !task done [task_ref] [desc]
``` 

The bot will then prompt you to define which task to update given a number, to give its number use
<!-- 
``` 
 !select [number]
```  -->

If you wish to avoid this step and just use the tasks id, do 

``` 
 !task done [task id] [desc]
``` 

This will then ping a lead to review said commands

