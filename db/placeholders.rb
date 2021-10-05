require "json"
require "csv"

# lookup rhe code in drupal data and extract summary / image info
def code_details(data, code)
   data.each do |rec|
      fund = rec['field_fund_id'].first['value']
      if fund.index(code) == 0
         summary = ""
         img = {width: 0, height: 0, title: "", alt: "", filename: ""}
         raw_sumary = rec['field_summary'].first
         if !raw_sumary.nil?
            summary = raw_sumary['value'].strip
         end
         raw_image_info = rec['field_bookplate_image'].first
         if !raw_image_info.nil?
            img_url = raw_image_info['url']
            filename = img_url.split("/").last
            img = {width: raw_image_info['width'], height: raw_image_info['height'],
                   title: raw_image_info['title'], alt: raw_image_info['alt'], filename: filename}
         end

         if summary == "" && img[:filename] == ""
            # puts "     No summary or image data for #{code} in drupal"
            return nil
         end

         return {code: code, summary: summary, image: img}
      end
   end
   # puts "     Nothing found for #{code} in drupal"
   return nil
end

file = File.open "drupal.json"
data = JSON.load file

# these are the mappings from full fund name to one or more codes
mappings = {}
csv = CSV.parse(File.read("fund-codes.csv"), headers: true)
csv.each do |row|
   if !mappings.key?(row[1])
      mappings[row[1]] = []
   end
   mappings[row[1]].push(row[0])
end

no_data = []
count = 0

# These are the full fund (collection) name from the solr index
File.readlines('solr-bookplates.txt').each do |line|
   count += 1
   clean = line.strip
   if !mappings.key?(clean)
      next
   end

   if mappings[clean].length == 1
      code = mappings[clean].first
      code_info = code_details(data, code)
      if code_info.nil?
         no_data << "#{clean}"
      else
         # puts "     MATCH: #{clean}:#{code} -> #{code_info}"
      end
      next
   end

   match = false
   mappings[clean].each do | code |
      code_info = code_details(data, code)
      if !code_info.nil?
         # puts "  MATCH: #{clean}:#{code} -> #{code_info}"
         # fund_codes.push( "('#{clean}', '#{code}')" )
         match = true
         # code_info[:title] = clean
         # collection_info << code_info
         break
      end
   end
   if !match
      # puts "     #{clean}: NO MAPPING FOUND (Multiple Options)"
      no_data << "#{clean}"
   end
end
puts "DONE: #{count} bookplate names."
# puts "NAMES WITH NO DATA:"
# puts no_data.join("\n")

id = 68
puts "INSERT into collections (id,filter_name,title) VALUES"
no_data.each do |ci|
   val = "   (#{id}, 'FilterFundCode', '#{ci}')"
   if ci == no_data.last
      puts "#{val};"
   else
      puts "#{val},"
   end
   id += 1
end


id = 68
puts
puts "INSERT into images (collection_id, alt_text, title, width, height, filename) VALUES"
no_data.each do |ci|
   val = "   (#{id}, 'University of Virginia Library', 'University of Virginia Library', 912, 1288, 'generic-bookplate.png' )"
   if ci == no_data.last
      puts "#{val};"
   else
      puts "#{val},"
   end
   id += 1
end

puts
puts "INSERT into images (collection_id, alt_text, title, width, height, filename) VALUES"
img_ids = [3, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 32, 33, 35, 36, 42, 44, 47, 48, 50, 52, 54, 55, 57, 59, 62, 63, 64, 66, 67]
for id in 2..67
   if img_ids.include?(id) == false
      puts "   (#{id}, 'University of Virginia Library', 'University of Virginia Library', 912, 1288, 'generic-bookplate.png' ),"
   end
end


