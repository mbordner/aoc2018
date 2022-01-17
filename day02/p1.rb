$twos = 0
$threes = 0
File.open("./data.txt").each_line do |line|
  if line.match(/^.*(.).*\1{1}.*\1{1}[^\1]*$/)
    $threes += 1
  end
  if line.match(//)
end