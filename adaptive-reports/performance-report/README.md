# Performance Report
This repository contains the code that generates the coaching report that Adaptive delivers at the beginning of each quarter. This report is designed to generate as many coaching and improvement ideas as possible.  Below are the various analyses and the list of possible responses generated for each analysis.

## Analysis Sentiment
Each recommendation is based on a numeric value that we then map to a form of sentiment that we then use as a color coding in the report -
1. Strong Red - Very concerning 
1. Weak Red - Somewhat concerning
1. Strong Yellow - Bordering on concern
1. Weak Yellow - Bordering on positive
1. Weak Green - Somewhat positive
1. Strong Green - Very positive
1. Neutral - No sentiment

## Analyses
Below are the various analyses that we built into the coaching report.  Note that the _%0.2f_ notation below will be replace by an actual number.
### Quantity
This analysis looks at the number of people that provided feedback.  A small number of people providing feedback is typically not healthy or desired. We want people to cultivate a healthy network of colleague advisors to help them level up their game.
#### Strong Red
3 pieces of feedback or fewer on average
##### Response Options
1. You don't have a lot of feedback. You received %.2f comments on average per topic. This quarter you should work with your coach to ensure you have more feedback at the end of this coming quarter.
1. It will be hard for you to make changes with limited feedback. You received %.2f comments on average per topic. Talk with your team about how you can improve.
1. You could use a lot more feedback. You received %.2f comments on average per topic. Talk with your team about what you can do to recieve more.
#### Weak Yellow
More than 3 pieces but less than 7
##### Response Options
1. Good job on your feedback. You received %.2f comments on average per topic. See if you can get more feedback next time.
1. Nice job on the amount of feedback you collected. You received %.2f comments on average per topic. I wonder if you can collect more feedback next time?
1. Your colleagues care about you! You received %.2f comments on average per topic. Come up with a way to collect more feedback next time
#### Weak Green
Between 7 and 9 pieces of feedback, inclusively 
##### Response Options
1. You received a lot of feedback!. You received %.2f comments on average per topic. Honor the time they spent by using this feedback to level up your game.?
1. Great job on your feedback!. You received %.2f comments on average per topic. See if there are themes that you can use for improving your performance.
1. Your colleagues gave you a lot of feedback. You received %.2f comments on average per topic. Show them you appreciate their feedback by making concrete improvements.
#### Strong Green
10 or more pieces of feedback
##### Response Options
1. Wow! You received a lot of feedback this quarter! You received %.2f comments on average per topic. Can you teach others how to do as well as you
1. Great job on your feedback. You received %.2f comments on average per topic. Can you use all this feedback to help the entire team?
1. Your colleagues clearly care about you.  You received %.2f comments on average per topic. Can you help others collect as much feedback?

## Overall Performance
This is the calculated aggregate score for all of the feedback for the quarter. The score tranches loosely follow grade school grades where 70's is a 'C', 80's is a 'B', 90's is an 'A' and above 95 is an 'A+'.  The Adaptive goal is to produce high B and A players.
#### Strong Red
Score less than 80
##### Response Options
1. You have a lot of room for improvement. Your score was %.2f.  Talk with your colleagues about where you can improve the most and come up with a plan to target those areas.
1. You can get your score way up. Your score was %.2f.  Think about where your greatest area of imptovement is and ask your colleagues for help.
1. You can do a lot better this coming quarter. Your score was %.2f. Find a coach that is really good in the areas where you can improve the most.
#### Weak Yellow
Between 80 and 90
##### Response Options
1. You have some of room for improvement. Your score was %.2f.  Talk with your colleagues about where you can improve the most and come up with a plan to target those areas.
1. You can get your score a bit. Your score was %.2f.  Think about where your greatest area of imptovement is and ask your colleagues for help.
1. You can do better this coming quarter. Your score was %.2f. Find a coach that is really good in the areas where you can improve the most.
#### Weak Green
Between 90 and 95
##### Response Options
1. Nice work this quarter! Your score was %.2f. Talk with your team about what you can do to get your scores even higher!
1. Good job! Your colleagues clearly have confidence in you. Your score was %.2f.  I wonder if you can find a way to do even better next time?
1. You turned in a solid quarter. Your score was %.2f. You colleagues have a lot of confidence in you. Think about how to use that confidence to turn in an even better performance next time.
#### StrongGreen
Greater than 95
##### Response Options
1. Wow! You really did great! Your score was %.2f.  You should think about how to stretch yourself. It looks like you may have reached the pinnacle of your performance
1. I'm impressed and I'm just a computer! Your score was %.2f.  I wonder if you can teach others how to earn the same level of confidence from your colleagues?
1. What an impressive job you turned in!. Your score was %.2f. How else can you use your skills to improve yourself and the company?

## Network Strength
This is a measure of the network that the person has cultivated that provided them with feedback.  A small network is typically a sign of a disconnected and disengaged employee.  A very large network can be a sign that someone is overwhelmed.
#### Strong Red
3 or smaller
##### Response Options
1. Your network was very limited of roughly %d people. Think about how you can expand your influence.
1. You don't appear to have many colleagues, only about %d people. Find a coach that can help you grow your network.
1. You could benefit from a larger network. Your network is roughly %d people.  I wonder how you can have a greater impact?
#### Weak Green
Between 4 and 6 inclusively
##### Response Options
1. You have a nice sized network of roughly %d people. You could probably expand your influence a little more.  Think about how you can do that this quarter.
1. Your network is a healthy size of roughly %d people.. I wonder if you can expand it just a little bit further?
1. You have a good sized network of roughly %d people.. Do you think you help a larger community be awesome?
#### StrongGreen
Between 7 and 9 inclusively
##### Response Options
1. You influence a lot of people, about %d people! That can be great but be careful not to overcommit yourself.
1. Your network is very healthy of roughly %d people.  Strengthen those relationships and be careful not to expand your network much more.
1. A lot of people care about you, about %d people! Make sure you return the favor by not overcommitting yourself.
#### Weak Yellow
Greater than 9
##### Response Options
1. Your network is very big, roughly %d people.  I wonder if you feel a bit overwhelmed? If so it is time to pair back and focus.",
1. A large network of people appear to rely on you, about %d people. You may want to consider focusing your energies.",
1. Your network is about %d people. Supporting too many people can water down your value to the company. If you feel overworked consider shedding lower value responsibilities.",

## Sentiment
Sentiment should generally be green.  If, even after being coached by the bot to write positive feedback, a person still recieves negative feedback, there is a serious problem afoot.
#### Strong Red
Very negative sentiment
##### Response Options
1. You have a lot of room for improvement. The sentiment of your feedback was %s. Work with a coach to come up with a plan.
1. You can do a lot better. The sentiment of your feedback was %s. Review your feedback and work with a coach to come up with a plan.
1. There is a lot you can do to change the sentiment of your feedback which was %s.  Review your feedback and think about what you can do to turn things around.
#### Weak Red
Negative sentiment
##### Response Options
1. You have some room for improvement. The sentiment of your feedback was %s. Work with a coach to come up with a plan.
1. You can do better. The sentiment of your feedback was %s. Review your feedback and work with a coach to come up with a plan.
1. There is a lot you can do to change the sentiment of your feedback which was %s.  Review your feedback and think about what you can do to turn things around.
#### Weak Yellow
Negative or neutral sentiment
##### Response Options
1. Your feedback was very factual. The sentiment of your feedback was %s. Talk with your colleagues about how to get more positive feedback.
1. I see that your feedback was very objective. The sentiment of your feedback was %s. Ask your colleagues what you can do to generate more positive feedback
1. I couldn't see any positive or negative feedback.  The sentiment of your feedback was %s. I wonder what you could do to generate more positie feedback?
#### Weak Green
Positive sentiment
##### Response Options
1. Nice job on your feedback! The sentiment of your feedback was %s. Now think about how to do a little better.
1. Well done! The sentiment of your feedback was %s. I wonder if you can do even better next time?
1. Nice work! The sentiment of your feedback was %s. I bet you can do even better next time.
#### StrongGreen
Very positive sentiment
##### Response Options
1. Wow! Your feedback was very positive. The sentiment of your feedback was %s. Can you teach others your super power?.
1. Very nice job! You received very positive feedback. The sentiment of your feedback was %s. What can you do to help the entire team get this kind of feedback
1. Very nice work! The sentiment of your feedback was %s. What can you do with this feedback to help the entire team?

## Energy
Energy is a measure of how much positive or negative energy you are producing for the team. The score tranches loosely follow grade school grades where 70's is a 'C', 80's is a 'B', 90's is an 'A' and above 95 is an 'A+'.  The Adaptive goal is to produce high B and A players.
#### Strong Red
Less than 80
##### Response Options
1. You could do a lot better creating positive energy for your colleagues. Your score was %.2f. Find a coach and brainstorm ideas.
1. You can improve the energy on your team a lot. Your score was %.2f. Think about what you can do to up your game.
1. You have a lot of room for improving the energy on your team. Your score was %.2f. Read your feedback and come up with a strategy.
#### Weak Yellow
Between 80 and 90
##### Response Options
1. You could do better creating positive energy for your colleagues. Your score was %.2f. Find a coach and brainstorm ideas.
1. You can improve the energy on your team a bit. Your score was %.2f. Think about what you can do to up your game.
1. You have some room for improving the energy on your team. Your score was %.2f. Read your feedback and come up with a strategy.
#### Weak Green
Between 90 and 95
##### Response Options
1. Nice job creating positive energy for your team. Your score was %.2f. A happy team is a more productive team! Think about how you can do even better.
1. Your colleagues like the energy you are generating for the team. Your score was %.2f. I wonder if you can do even better?
1. You are helping to create a hapy team. Your score was %.2f. Nice work! Now come up with a plan to improve your energy a little more.
#### StrongGreen
Above 95
##### Response Options
1. Wow! You really generated a lot of great energy! Your score was %.2f. What can you do to teach others your super power?
1. Your team loves you! Your score was %.2f. I wonder if you can leverage your positive energy to create a sustainble change in the team?
1. Goodness you created a lot of positive energy. Your score was %.2f. Do you think you can help the entire team do what you do?
## Improvement
The perception of imporvement helps everyone appreciate the effort that someone is putting into their growth. Helping your colleagues know what you are working on can greatly influence this score and create a virtuous cycle of improvement.
#### Strong Red
Got a lot worse
##### Response Options
			"Your colleagues feel like you could improve a lot.  Think about what you will do differently next time",
			"It looks like you have a lot of room for improvement.  Talk with your team about what you can do to turn things around.",
			"You can do a lot better to improve.  Read your feedback and come up with a plan.",
#### Weak Yellow
Got worse
##### Response Options
1. Your colleagues feel like you could improve a bit.  Think about what you will do differently next time.
1. It looks like you have some room to improve.  Talk with your team about what you can do to turn things around.
1. You can do better next time to improve.  Read your feedback and come up with a plan.
#### Weak Green
Got better
##### Response Options
1. Nice work. Your colleagues feel like you are improving! What can you do to do even better?
1. It looks like your colleagues feel like you are improving. I wonder what you can do to improve that impression?
1. Your colleagues see your improvements.  Good job. Do you think you can do even better next time?
#### StrongGreen
Got a lot better
##### Response Options
1. Great work! Your colleagues feel like you improved as lot! Do you think you can help others improve like you did?
1. It looks like your colleagues feel like you are improving a lot! Can you teach others your super power?
1. Yor colleagues see big improvements!  Great job. Do you think you help the entire team improve like you did?

## Consistency
Consistency is a measure of the variance in performance quarter to quarter.  This is perhaps the most powerful metric for identifying struggling employees. When the variance is large is indicates someone whose performance is great one quarter and terrible the next.  If their variance is very low and their scores are bad then you probably want to terminate their employment.  If their variance is low and their scores are great then you have a bedrock employee that you want to bring more of magic to the company.
#### Strong Red
Greater than 80
##### Response Options
1. Your performance is rather inconsistent. This can negatively impact your team.  Think about ways to stabilize your performance.",
1. Your colleagues would benefit from more stable performance quarter to quarter. I wonder how you can turn in a more consisent performance?",
1. It looks like your performance is up and down quarter to quarter. Try to find ways of improving your consistency.",
#### Weak Yellow
Between 60 and 80
##### Response Options
1. Your performance is relatively stable and I think you can do better.  Think about ways to improve the conistency of your performance.
1. Your performance is consistent quarter to quarter. Try to improve that consistency next time.
1. Performing consistently contributes to stable teams and you did a good job.  I wonder if you can do better next time.
#### Weak Green
Between 40 and 60
##### Response Options
1. Nice work! Your colleagues see you as a consistent performer! What can you do to do even better?
1. It looks like your colleagues see you as a consistent performer. I wonder what you can do to improve that impression?
1. Your colleagues like your consistency.  Good job. Do you think you can do even better next time?
#### StrongGreen
Less than 40
##### Response Options
1. Great work! Your colleagues see you as a consistent performer! Do you think you can help others be as consistent as you?
1. It looks like your colleagues feel like you are a very consistent performer! Can you teach others your super power?
1. Your colleagues recognize your consistency!  Great job. Do you think you help the entire team achieve the same consistency?
