#!/usr/bin/env coffee
 
 
fs = require 'fs'
path = require 'path'
 
 
copyEntries = (src, dst) ->
  fs.readdir src, (err, files) ->
    for file in files
      s = path.join src, file
      d = path.join dst, file
      copy s, d
 
 
copy = (src, dst) ->
  fs.lstat src, (err, stat) ->
    switch true
      when stat.isFile()
        fs.link src, dst, (err) ->
          if err? and err.code != 'EEXIST'
            console.log err.stack
 
      when stat.isDirectory()
        fs.mkdir dst, 0o0755, (err) ->
          if err? and err.code != 'EEXIST'
            console.log err.stack
          else
            copyEntries src, dst
 
      when stat.isSymbolicLink()
        fs.readlink src, (err, lsrc) ->
          if err? and err.code != 'EEXIST'
            console.log err.stack
          else
            fs.symlink lsrc, dst
 
      else
        console.log "#{src} not handled (#{stat})"
    return null # tiny optimization
 
 
if require.main is module
  args = process.argv[2..]
  if args.length != 2
    console.log "invalid usage!"
    process.exit 1
  else
    fs.mkdir args[1], 0o0755, (err) ->
      if err and err.code != 'EEXIST'
        console.log err.stack
      else
        copy args...
