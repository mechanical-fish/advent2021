#!/usr/bin/env ruby

x = 0
y = 0
File.readlines('input.txt').each do |d|
  p = d.split(' ')
  x += p[1].to_i if p[0] == 'forward'
  y += p[1].to_i if p[0] == 'down'
  y -= p[1].to_i if p[0] == 'up'
end
puts "Distance is #{x}, depth is #{y}, product is #{x*y}"

aim = 0
depth = 0
pos = 0
File.readlines('input.txt').each do |d|
  p = d.split(' ')
  n = p[1].to_i
  if p[0] == 'forward'
    pos += n
    depth += aim * n
  elsif p[0] == 'down'
    aim += n
  elsif p[0] == 'up'
    aim -= n
  end
end
puts "Distance is #{pos}, depth is #{depth}, product is #{pos*depth}"
