using famousfront.utils;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Diagnostics.CodeAnalysis;
using System.Diagnostics.Contracts;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Media;
namespace famousfront.utils
{
  public static class ElementOptions
  {
    #region Font size

    public static readonly DependencyProperty TitleFontSizeProperty =
        DependencyProperty.RegisterAttached("TitleFontSize", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(12d * (96d / 72d),
                                                                          FrameworkPropertyMetadataOptions.AffectsMeasure |
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetTitleFontSize(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return (double)obj.GetValue(TitleFontSizeProperty);
    }

    public static void SetTitleFontSize(DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(TitleFontSizeProperty, value);
    }

    public static readonly DependencyProperty HeaderFontSizeProperty =
        DependencyProperty.RegisterAttached("HeaderFontSize", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(13d * (96d / 72d),
                                                                          FrameworkPropertyMetadataOptions.AffectsMeasure |
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetHeaderFontSize( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(HeaderFontSizeProperty));
    }

    public static void SetHeaderFontSize( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(HeaderFontSizeProperty, value);
    }

    public static readonly DependencyProperty ContentFontSizeProperty =
        DependencyProperty.RegisterAttached("ContentFontSize", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(11d * (96d / 72d),
                                                                          FrameworkPropertyMetadataOptions.AffectsMeasure |
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetContentFontSize( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(ContentFontSizeProperty));
    }

    public static void SetContentFontSize( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ContentFontSizeProperty, value);
    }

    public static readonly DependencyProperty TextFontSizeProperty =
        DependencyProperty.RegisterAttached("TextFontSize", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(10d * (96d / 72d),
                                                                          FrameworkPropertyMetadataOptions.AffectsMeasure |
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetTextFontSize( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(TextFontSizeProperty));
    }

    public static void SetTextFontSize( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(TextFontSizeProperty, value);
    }

    #endregion

    #region Thickness

    public static readonly DependencyProperty DefaultThicknessProperty =
        DependencyProperty.RegisterAttached("DefaultThickness", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(1d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetDefaultThickness( DependencyObject obj)
    {
      Debug.Assert(obj != null);
//      ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(DefaultThicknessProperty));
    }

    public static void SetDefaultThickness( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(DefaultThicknessProperty, value);
    }

    public static readonly DependencyProperty SemiBoldThicknessProperty =
        DependencyProperty.RegisterAttached("SemiBoldThickness", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(1.5d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetSemiBoldThickness( DependencyObject obj)
    {
      Debug.Assert(obj != null);
 //     ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(SemiBoldThicknessProperty));
    }

    public static void SetSemiBoldThickness( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(SemiBoldThicknessProperty, value);
    }

    public static readonly DependencyProperty BoldThicknessProperty =
        DependencyProperty.RegisterAttached("BoldThickness", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(2d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetBoldThickness( DependencyObject obj)
    {
      Debug.Assert(obj != null);
//      ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(BoldThicknessProperty));
    }

    public static void SetBoldThickness( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(BoldThicknessProperty, value);
    }

    public static readonly DependencyProperty DefaultThicknessValueProperty =
        DependencyProperty.RegisterAttached("DefaultThicknessValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(1d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetDefaultThicknessValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(DefaultThicknessValueProperty));
    }

    public static void SetDefaultThicknessValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(DefaultThicknessValueProperty, value);
    }

    public static readonly DependencyProperty SemiBoldThicknessValueProperty =
        DependencyProperty.RegisterAttached("SemiBoldThicknessValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(1.5d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetSemiBoldThicknessValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(SemiBoldThicknessValueProperty));
    }

    public static void SetSemiBoldThicknessValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(SemiBoldThicknessValueProperty, value);
    }

    public static readonly DependencyProperty BoldThicknessValueProperty =
        DependencyProperty.RegisterAttached("BoldThicknessValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(2d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetBoldThicknessValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(BoldThicknessValueProperty));
    }

    public static void SetBoldThicknessValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(BoldThicknessValueProperty, value);
    }

    #endregion

    #region Padding

    public static readonly DependencyProperty DefaultPaddingProperty =
        DependencyProperty.RegisterAttached("DefaultPadding", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(1d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetDefaultPadding( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(DefaultPaddingProperty));
    }

    
    public static void SetDefaultPadding( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(DefaultPaddingProperty, value);
    }

    
    public static readonly DependencyProperty SemiBoldPaddingProperty =
        DependencyProperty.RegisterAttached("SemiBoldPadding", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(4d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    
    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetSemiBoldPadding( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(SemiBoldPaddingProperty));
    }

    
    public static void SetSemiBoldPadding( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(SemiBoldPaddingProperty, value);
    }

    
    public static readonly DependencyProperty BoldPaddingProperty =
        DependencyProperty.RegisterAttached("BoldPadding", typeof(Thickness), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Thickness(8d),
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, ThicknessUtil.CoerceNonNegative));

    
    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static Thickness GetBoldPadding( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      ThicknessUtil.EnsureNonNegative();
      return BoxingHelper<Thickness>.Unbox(obj.GetValue(BoldPaddingProperty));
    }

    
    public static void SetBoldPadding( DependencyObject obj, Thickness value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(BoldPaddingProperty, value);
    }

    
    public static readonly DependencyProperty DefaultPaddingValueProperty =
        DependencyProperty.RegisterAttached("DefaultPaddingValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(1d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    
    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetDefaultPaddingValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(DefaultPaddingValueProperty));
    }

    
    public static void SetDefaultPaddingValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(DefaultPaddingValueProperty, value);
    }

    
    public static readonly DependencyProperty SemiBoldPaddingValueProperty =
        DependencyProperty.RegisterAttached("SemiBoldPaddingValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(4d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    
    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetSemiBoldPaddingValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(SemiBoldPaddingValueProperty));
    }

    
    public static void SetSemiBoldPaddingValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(SemiBoldPaddingValueProperty, value);
    }

    
    public static readonly DependencyProperty BoldPaddingValueProperty =
        DependencyProperty.RegisterAttached("BoldPaddingValue", typeof(double), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(8d,
                                                                          FrameworkPropertyMetadataOptions.AffectsArrange |
                                                                          FrameworkPropertyMetadataOptions.Inherits,
                                                                          null, DoubleExtension.CoerceNonNegative));

    
    [SuppressMessage("Microsoft.Contracts", "Ensures", Justification = "Can't be proven.")]
    public static double GetBoldPaddingValue( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      DoubleExtension.EnsureNonNegative();
      return BoxingHelper<double>.Unbox(obj.GetValue(BoldPaddingValueProperty));
    }

    
    public static void SetBoldPaddingValue( DependencyObject obj, double value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(BoldPaddingValueProperty, value);
    }

    #endregion

    #region Animation

    
    public static readonly DependencyProperty DefaultDurationProperty =
        DependencyProperty.RegisterAttached("DefaultDuration", typeof(Duration), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Duration(TimeSpan.FromSeconds(0d)), FrameworkPropertyMetadataOptions.Inherits));

    
    public static Duration GetDefaultDuration( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return BoxingHelper<Duration>.Unbox(obj.GetValue(DefaultDurationProperty));
    }

    
    public static void SetDefaultDuration( DependencyObject obj, Duration value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(DefaultDurationProperty, value);
    }

    
    public static readonly DependencyProperty MinimumDurationProperty =
        DependencyProperty.RegisterAttached("MinimumDuration", typeof(Duration), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Duration(TimeSpan.FromSeconds(0.2d)),
                                                                          FrameworkPropertyMetadataOptions.Inherits));

    
    public static Duration GetMinimumDuration( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return BoxingHelper<Duration>.Unbox(obj.GetValue(MinimumDurationProperty));
    }

    
    public static void SetMinimumDuration( DependencyObject obj, Duration value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(MinimumDurationProperty, value);
    }

    
    public static readonly DependencyProperty OptimumDurationProperty =
        DependencyProperty.RegisterAttached("OptimumDuration", typeof(Duration), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Duration(TimeSpan.FromSeconds(0.5d)),
                                                                          FrameworkPropertyMetadataOptions.Inherits));

    
    public static Duration GetOptimumDuration( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return BoxingHelper<Duration>.Unbox(obj.GetValue(OptimumDurationProperty));
    }

    
    public static void SetOptimumDuration( DependencyObject obj, Duration value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(OptimumDurationProperty, value);
    }

    
    public static readonly DependencyProperty MaximumDurationProperty =
        DependencyProperty.RegisterAttached("MaximumDuration", typeof(Duration), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new Duration(TimeSpan.FromSeconds(1.0d)),
                                                                          FrameworkPropertyMetadataOptions.Inherits));

    
    public static Duration GetMaximumDuration( DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return BoxingHelper<Duration>.Unbox(obj.GetValue(MaximumDurationProperty));
    }

    
    public static void SetMaximumDuration( DependencyObject obj, Duration value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(MaximumDurationProperty, value);
    }

    #endregion

    #region Render Resources
    public static readonly DependencyProperty ContrastBrushProperty =
        DependencyProperty.RegisterAttached("ContrastBrush", typeof(SolidColorBrush), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new SolidColorBrush(Color.FromArgb(0xff, 0xEE, 0xEE, 0xEE)),
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits));


    public static SolidColorBrush GetContrastBrush(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return (SolidColorBrush)obj.GetValue(ContrastBrushProperty);
    }


    public static void SetContrastBrush(DependencyObject obj, SolidColorBrush value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ContrastBrushProperty, value);
    }

    public static readonly DependencyProperty ForegroundBrushProperty =
        DependencyProperty.RegisterAttached("ForegroundBrush", typeof(SolidColorBrush), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new SolidColorBrush(Color.FromArgb(0xff, 0x11, 0x11, 0x11)),
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits));


    public static SolidColorBrush GetForegroundBrush(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return (SolidColorBrush)obj.GetValue(ForegroundBrushProperty);
    }


    public static void SetForegroundBrush(DependencyObject obj, SolidColorBrush value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ForegroundBrushProperty, value);
    }
    
    public static readonly DependencyProperty BackgroundBrushProperty =
        DependencyProperty.RegisterAttached("BackgroundBrush", typeof(SolidColorBrush), typeof(ElementOptions),
                                            new FrameworkPropertyMetadata(new SolidColorBrush(Color.FromArgb(0xff, 0xEE, 0xEE, 0xEE)),
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits));


    public static SolidColorBrush GetBackgroundBrush(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return (SolidColorBrush)obj.GetValue(BackgroundBrushProperty);
    }


    public static void SetBackgroundBrush(DependencyObject obj, SolidColorBrush value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(BackgroundBrushProperty, value);
    }

    #endregion
  }
  [DebuggerNonUserCode]
  [System.Diagnostics.Contracts.Pure]
  internal static class BoxingHelper<T>
      where T : struct
  {
    internal static T Unbox(object value)
    {
      Contract.Assume(value is T);
      return (T)value;
    }
  }
}
