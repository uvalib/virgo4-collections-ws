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

missed = []
fund_codes = []
collection_info = []
count = 0
# These are the full fund (collection) name from the solr index
File.readlines('solr-bookplates.txt').each do |line|
   count += 1
   clean = line.strip
   # puts "=====> SOLR COLLECTION [#{clean}]"
   if !mappings.key?(clean)
      missed << "#{clean} (no code)"
      next
   end

   if mappings[clean].length == 1
      code = mappings[clean].first
      code_info = code_details(data, code)
      if code_info.nil?
         missed << "#{clean} : #{code}"
      else
         # puts "     MATCH: #{clean}:#{code} -> #{code_info}"
         fund_codes.push( "('#{clean}', '#{code}')" )
         code_info[:title] = clean
         collection_info << code_info
      end
      next
   end

   match = false
   # puts "  MULTIPLE CODES FOR [#{clean}]"
   mappings[clean].each do | code |
      code_info = code_details(data, code)
      if !code_info.nil?
         # puts "  MATCH: #{clean}:#{code} -> #{code_info}"
         fund_codes.push( "('#{clean}', '#{code}')" )
         match = true
         code_info[:title] = clean
         collection_info << code_info
         break
      end
   end
   if !match
      # puts "     #{clean}: NO MAPPING FOUND (Multiple Options)"
      missed << "#{clean} : #{ mappings[clean].join(", ")}"
   end
end
puts "DONE: #{count} bookplate names. MISSING: #{missed.length}, FOUND: #{fund_codes.length}"
# puts "MISSING NAMES:"
# puts missed.join("\n")
# puts FUND CODES:"
# puts fund_codes.join("\n")


# fc_vals = fund_codes.join(",\n")
# puts "INSERT into fund_codes (name,fund_code) VALUES #{fc_vals};"

# id = 2
# puts "INSERT into collections (id,description,filter_name,title) VALUES"
# collection_info.each do |ci|
#    summary = ci[:summary].gsub(/\'/, "''")
#    val = "   (#{id}, '#{summary}', 'FilterFundCode', '#{ci[:code]}')"
#    if ci == collection_info.last
#       puts "#{val};"
#    else
#       puts "#{val},"
#    end
#    id += 1
# end


# id = 2
# puts "INSERT into images (collection_id, alt_text, title, width, height, filename) VALUES"
# collection_info.each do |ci|
#    img = ci[:image]
#    if img[:filename] != ""
#       val = "   (#{id}, '#{img[:alt]}', '#{img[:title]}', #{img[:width]}, #{img[:height]}, '#{img[:filename]}' )"
#       if ci == collection_info.last
#          puts "#{val};"
#       else
#          puts "#{val},"
#       end
#    end
#    id += 1
# end

# id = 2
# puts "INSERT into collection_features (collection_id, feature_id) VALUES"
# collection_info.each do |ci|
#    val = "   ( #{id}, 5 )"
#    if ci == collection_info.last
#       puts "#{val};"
#    else
#       puts "#{val},"
#    end
#    id += 1
# end

id = 2
collection_info.each do |ci|
   title = ci[:title].gsub(/\'/, "''")
   puts "UPDATE collections set title='#{title}' where id=#{id};"
   id += 1
end

