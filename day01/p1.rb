$freq = 0
File.read("./data.txt").each_line do |line|
  $freq += line.to_i
end
puts $freq