# Algorithms

The requirement consists of a min and max budget or sometimes just either of them.

BUDGET CALC 

>> Min and max budget given <<

  min = Max( (min_b - (25% of min_b), 1.0)
                 
  max = Max((max_b + (25% of max_b)), 1.25)

where min_b= minimum budget
      max_b= maximum budget
>> only min/max <<

  >> only min,
  min = max((min_b - (25% of min_b)), 1.0)
  max = max( (min_b + (25% of min_b)), 1.25)

  >>only max,
  min = 0 
  max = max( (max_b + (25% of max_b)), 1.25)
 
If the budget is within min and max budget, budget contribution for the match percentage is full 30%. If min or max is not given, +/- 10% budget is a full 30% match.
In order to achieve full 30% match, we should distribute the 30 marks in the range min-minimum budget, and max budget - max, where min_b and max_b gets full 30 marks and min and max gets 0.


% x = ((x - min) / (min_b-min)) * 100

% x = ((max - x) / (max-max_b)) * 100


Any price falling within the +-10% of minBudget or maxBudget given in the requirement will get 30 marks.

% x = ((x - min) / (min_b10-min)) * 100

% x = ((max - x) / (max-max_b10)) * 100

where min_b10= price falling within the +-10% of min budget
      max_b10=price falling within the +-10% of max budget


BEDROOM FILTERING

The requirement consists of a min and max bedrooms or sometimes just either of them. We need to consider only those property listings which are within +-2 of the requirement bedrooms.

>> both min and max bedrooms <<

Then range of bedrooms attribute in the properties table that should be matched can be calculated by

  min_bedrooms	= max((min_bedroomsRequired - 2), 1)  
                   
  max_bedrooms	= max((max_bedroomsRequired + 2), 1)

>> only min/max bedrooms <<
   >> only min
  min_bedrooms	= max((min_bedroomsRequired - 2), 1)
  max_bedrooms	= max((min_bedroomsRequired + 2), 1)

  >> only max
  min_bedrooms	= max((max_bedroomsRequired - 2), 1)
  max_bedrooms	= max((max_bedroomsRequired + 2), 1)

If bedroom and bathroom fall between min and max, each will contribute full 20%. If min or max is not given, match percentage varies according to the value.

Similar to the budget calc, the bedroom/bathroom is also calculated according to 20%.  20 marks in the range minBedrooms - min bedrooms required, and max bedrooms required - max bedrooms, where min bedrooms required and max bedrooms required gets full 20 marks and min bedrooms-1 and max bedrooms+1 get 0.

For range minBedrooms - minBedroomsRequired

 % x = ( (x - min_bedrooms) / (min_bedroomsRequired - min_bedrooms) ) * 100

  % x = ( (max_bedrooms - x) / (max_bedrooms - max_bedroomsRequired) ) * 100

BATHROOM FILTERING

It is exactly the same as bedrooms filtering


>> both min and max bathrooms <<

Then range of bathrooms attribute in the properties table that should be matched can be calculated by

  min_rooms	= max((min_roomsRequired - 2), 1)  
                   
  max_rooms	= max((max_roomsRequired + 2), 1)

where max_rooms= Maximum bathrooms
      min_rooms=Minimum bathrooms

>> only min/max bathrooms <<
   >> only min
  min_rooms	= max((min_roomsRequired - 2), 1)
  max_rooms	= max((min_roomsRequired + 2), 1)


  >> only max
  min_rooms	= max((max_roomsRequired - 2), 1)
  max_rooms	= max((max_roomsRequired + 2), 1)

If bedroom and bathroom fall between min and max, each will contribute full 20%. If min or max is not given, match percentage varies according to the value


For range min bathrooms - min bathroomsRequired

  % x = ( (x - min_rooms) / (min_roomsRequired - min_rooms) ) * 100

  % x = ( (max_rooms - x) / (max_rooms - max_roomsRequired) ) * 100

Finally the distance,

DISTANCE CALC

The calc of distance is different than the max and min buget or bedrooms.
According to mathematics, the formula to calc the distance is as follows,

 D = acos( sin(lat_1)*sin(lat_2)
 + cos(lat_1)*cos(lat_2)*cos(lon_2 - lon_1) ) * R
Where,

 lat_1 & lat_2   =  latitude in radians,
 lon_1 & lon_2 	=  longitude in radians,
 R 	        =  radius of earth in miles
 D 	        =  distance in miles

We can make a bounding box which is depicted by the blue square in the above figure. We can get those bounding latitude and longitude by the below mentioned formula:

 Lat_top 	= 	lat_given + (d/R)
 Lat_bottom      = 	lat_given - (d/R)
 Lon_left 	= 	lon_given - asin(d/R)/cos(lat_given)
 Lon_right 	= 	lon_given + asin(d/R)/cos(lat_given)


so the final SQL query will be 

 SELECT latitude, longitude, acos( sin(latitude)*sin(lat_given)  +  
 cos(latitude)*cos(lat_given)*cos(lon_given - longitude) ) * R 
 as distance, price, bedrooms, bathrooms
 From properties
 Where latitude Between lat_top And lat_bottom
 And longitude Between lon_left And lon_right
 And distance <= 10
 And price Between min And max
 And bedrooms Between min_bedrooms And max_bedrooms
 And bathrooms Between minrooms And maxrooms;

After giving the filtering criteria, we can write a sorted query with all the 4 parameters in it, since it takes lesser time to give the results than doing it seperately.
Basically, the first step includes creatind a DB including all the 4 parameters in it. Followed by creating the matching algorithm.
The matching Algorithm can be carried over through Go routines Or the standard matching algorithm.
The matching method and its corresponding matching algorithms are part of the matching criteria. They help to determine how a specific field in one record is compared to the same field in another record and whether the fields are considered matches.
When the primary matching criteria is done. Then we can match according to the weightage like the distance and budget criteria.
As part of the process of creating match key values, matching rule field values are normalized. How a field value is normalized depends on several factors, including the matching method for that field as specified in the matching rule. In addition, some commonly used fields are normalized to optimize duplicate detection.
