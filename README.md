# Discutils

A package that provides some nice utils for discordgo.

## Events

This package provides some wrappers for events. Several places can register a handler. The advantage of this over the integrated solution is that these handlers can prevent further handlers from starting.

## Commands

Commands are handled a bit differentially - they are parsed into an argument -> value map before being passed down to the command handler. The same is done for modals, autocomplete, etc

## Interactions

There are abstractions for interactions - just some utils I find useful to work with. Check docs for those.

## Getters

There are getters for channels, messages and users - the idea being to fetch from cache (state) first, and if that fails, fetch from the api.