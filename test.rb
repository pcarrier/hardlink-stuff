#!/usr/bin/env ruby
 
require 'find'
 
exit 1 unless ARGV.length == 2
 
SRC, DST = ARGV
 
def clone_dir srcdir, dstdir
  Dir.entries(srcdir).each do |e|
    next if %w[. ..].include? e
    src = File.join srcdir, e
    dst = File.join dstdir, e
    stat = File::lstat src
    case true
    when stat.directory?
      Dir.mkdir dst
      clone_dir src, dst
    when stat.file?
      File.link src, dst
    when stat.symlink?
      File.symlink File.readlink(src), dst
    end
  end
end
 
Dir.mkdir DST
clone_dir SRC, DST
