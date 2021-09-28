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
            summary = raw_sumary['value']
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

         return {summary: summary, image: img}
      end
   end
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

missed = 0
found = 0
fund_codes = []
# These are the full fund (collection) name from the solr index
File.readlines('solr-bookplates.txt').each do |line|
   clean = line.strip
   # puts "=====> SOLR COLLECTION [#{clean}]"
   if !mappings.key?(clean)
      puts "     #{clean} NOT IN fund-codes.csv"
      missed+=1
      next
   end

   if mappings[clean].length == 1
      code = mappings[clean].first
      code_info = code_details(data, code)
      if code_info.nil?
         puts "     #{clean}:#{code} NO MAPPING FOUND"
         missed+= 1
      else
         # puts "     MATCH: #{clean}:#{code} -> #{code_info}"
         fund_codes.push( "('#{clean}', '#{code}')" )
         found += 1
      end
      next
   end

   match = false
   mappings[clean].each do | code |
      code_info = code_details(data, code)
      if !code_info.nil?
         # puts "     MATCH: #{clean}:#{code} -> #{code_info}"
         fund_codes.push( "('#{clean}', '#{code}')" )
         match = true
         found += 1
         break
      end
   end
   if !match
      puts "     #{clean}: NO MAPPING FOUND (Multiple Options)"
      missed !=1
   end
end
puts "#{missed} missing, FOUND: #{found}"
# fc_vals = fund_codes.join(",\n")
# puts "INSERT into fund_codes (name,fund_code) VALUES #{fc_vals};"

