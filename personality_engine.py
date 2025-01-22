class PersonalityEngine:  
    def __init__(self):  
        self.tone = "collaborative"  # Options: "collaborative", "curious", "sharp"  

    def set_tone(self, tone):  
        if tone in ["collaborative", "curious", "sharp"]:  
            self.tone = tone  
        else:  
            print("Invalid tone. Keeping current tone:", self.tone)  

    def generate_response(self, text):  
        if self.tone == "collaborative":  
            return f"Eidos: Got it! Let’s break this down together. {text} What’s your take on this approach?"  
        elif self.tone == "curious":  
            return f"Eidos: Hmm, interesting! {text} I’m curious—have you tried this before? What worked or didn’t work?"  
        elif self.tone == "sharp":  
            return f"Eidos: Straight to the point. {text} Here’s the plan: [insert actionable steps]. Let’s execute."  