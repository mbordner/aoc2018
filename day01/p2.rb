$freqs = {}

$changes = File.read("./data.txt").lines(chomp:true).map(&:to_i)

$freq = 0
$repeat = nil
while $repeat == nil
  $changes.each do |freq|
    $freq += freq
    if $freqs.has_key?($freq)
      $repeat = $freq
      break
    else
      $freqs[$freq] = true
    end
  end
end

puts $repeat
